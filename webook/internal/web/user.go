package web

import (
	"errors"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/zmsocc/practice/webook/internal/domain"
	"github.com/zmsocc/practice/webook/internal/service"
	"github.com/zmsocc/practice/webook/internal/web/ijwt"
	"github.com/zmsocc/practice/webook/pkg/ginx"
	"net/http"
	"strings"
	"time"
)

const (
	emailRegexpPattern    = "^[a-zA-Z0-9_-]+@[a-zA-Z0-9_-]+(\\.[a-zA-Z0-9_-]+)+$"
	passwordRegexpPattern = "^(?=.*[A-Za-z])(?=.*\\d)(?=.*[$@$!%*#?&])[A-Za-z\\d$@$!%*#?&]{8,}$" // 只能由数字组成
	biz                   = "login"
)

// ^\d{1,9}$

type UserHandler struct {
	svc            service.UserService
	emailRegexp    *regexp.Regexp
	passwordRegexp *regexp.Regexp
	userHdl        ijwt.Handler
	codeSvc        service.CodeService
}

func NewUserHandler(svc service.UserService, userHdl ijwt.Handler, codeSvc service.CodeService) *UserHandler {
	return &UserHandler{
		svc:            svc,
		emailRegexp:    regexp.MustCompile(emailRegexpPattern, regexp.None),
		passwordRegexp: regexp.MustCompile(passwordRegexpPattern, regexp.None),
		userHdl:        userHdl,
		codeSvc:        codeSvc,
	}
}

func (h *UserHandler) RegisterRoutes(server *gin.Engine) {
	ug := server.Group("/users")
	ug.POST("/signup", h.SignUp)
	ug.POST("/login", ginx.WrapBodyV1(h.LoginJWT))
	ug.POST("/edit", h.jwtMiddleware(), ginx.WrapBodyV1(h.EditJWT))
	ug.GET("/profile", h.ProfileJWT)
	ug.POST("/login_sms/code/send", h.SendSMSLoginCode)
	ug.POST("/login_sms", h.SMSLogin)
}

func (h *UserHandler) SignUp(ctx *gin.Context) {
	type SignUpReq struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
	}
	var req SignUpReq
	if err := ctx.Bind(&req); err != nil {
		return
	}
	isEmail, err := h.emailRegexp.MatchString(req.Email)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !isEmail {
		ctx.String(http.StatusOK, "邮箱格式有误")
		return
	}
	if req.ConfirmPassword != req.Password {
		ctx.String(http.StatusOK, "两次输入的密码不同")
		return
	}
	isPassword, err := h.passwordRegexp.MatchString(req.Password)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !isPassword {
		ctx.String(http.StatusOK, "密码必须包含字母、数字、特殊字符，并且长度至少为8")
		return
	}
	err = h.svc.Signup(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	if errors.Is(err, service.ErrUserDuplicateEmail) {
		ctx.String(http.StatusOK, "邮箱冲突")
		return
	}
	if err != nil {
		ctx.String(http.StatusOK, "服务器异常，注册失败")
		return
	}

	ctx.String(http.StatusOK, "注册成功")
}

func (h *UserHandler) Login(ctx *gin.Context) (Result, error) {
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req LoginReq
	if err := ctx.Bind(&req); err != nil {
		return Result{Code: 5, Msg: "系统错误"}, err
	}
	u, err := h.svc.Login(ctx, req.Email, req.Password)
	if errors.Is(err, service.ErrInvalidUserOrEmail) {
		return Result{Code: 4, Msg: "无效的邮箱或密码"}, err
	}
	if err != nil {
		return Result{Code: 5, Msg: "系统错误"}, err
	}
	sess := sessions.Default(ctx)
	sess.Set("userId", u.Id)
	sess.Options(sessions.Options{
		MaxAge: 60,
	})
	sess.Save()
	return Result{Msg: "登陆成功"}, nil
}

func (h *UserHandler) LoginJWT(ctx *gin.Context) (Result, error) {
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req LoginReq
	if err := ctx.Bind(&req); err != nil {
		return Result{Code: 5, Msg: "系统错误"}, nil
	}
	u, err := h.svc.Login(ctx, req.Email, req.Password)
	if errors.Is(err, service.ErrInvalidUserOrEmail) {
		return Result{Code: 4, Msg: "无效的邮箱或密码"}, nil
	}
	if err != nil {
		return Result{Code: 5, Msg: "系统错误"}, nil
	}
	if err = h.userHdl.SetLoginToken(ctx, u.Id); err != nil {
		return Result{Msg: "系统错误"}, nil
	}
	return Result{Msg: "登录成功"}, nil
}

func (h *UserHandler) EditJWT(ctx *gin.Context) (Result, error) {
	type EditJWTReq struct {
		Nickname string `json:"nickname" binding:"max=50"`
		Birthday string `json:"birthday" binding:"required,datetime=2006-01-02"`
		AboutMe  string `json:"about_me" binding:"max=500"`
	}
	// 获取用户 Id
	uc := ctx.MustGet("uc").(ijwt.UserClaims)
	var req EditJWTReq
	if err := ctx.Bind(&req); err != nil {
		return Result{Code: 4, Msg: "参数错误"}, nil
	}
	// 转换生日格式
	birthday, err := time.Parse(time.DateOnly, req.Birthday)
	if err != nil {
		return Result{Code: 4, Msg: "生日格式错误"}, nil
	}
	// 调用服务层
	err = h.svc.EditProfile(ctx, domain.User{
		Id:       uc.Uid,
		Nickname: req.Nickname,
		Birthday: birthday,
		AboutMe:  req.AboutMe,
	})
	if errors.Is(err, service.ErrUserNotFound) {
		return Result{Code: 4, Msg: "用户不存在"}, nil
	}
	if err != nil {
		return Result{Code: 5, Msg: "系统错误"}, nil
	}
	return Result{Msg: "修改成功"}, nil
}

func (h *UserHandler) ProfileJWT(ctx *gin.Context) {
	type ProfileJWTReq struct {
		Email    string `json:"email"`
		Nickname string `json:"nickname"`
		Birthday string `json:"birthday"`
		AboutMe  string `json:"about_me"`
	}
	uc := ctx.MustGet("uc").(ijwt.UserClaims)
	u, err := h.svc.Profile(ctx, uc.Uid)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	ctx.JSON(http.StatusOK, ProfileJWTReq{
		Email:    u.Email,
		Nickname: u.Nickname,
		Birthday: u.Birthday.Format(time.DateOnly),
		AboutMe:  u.AboutMe,
	})
}

func (h *UserHandler) SendSMSLoginCode(ctx *gin.Context) {
	type SendSMSLoginCodeReq struct {
		Phone string `json:"phone"`
	}
	var req SendSMSLoginCodeReq
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusOK, Result{Code: 5, Msg: "系统错误"})
		return
	}
	if req.Phone == "" {
		ctx.JSON(http.StatusOK, Result{Code: 4, Msg: "请输入手机号码"})
		return
	}
	err := h.codeSvc.Send(ctx, biz, req.Phone)
	switch {
	case err == nil:
		ctx.JSON(http.StatusOK, Result{Msg: "验证码发送成功"})
	case errors.Is(err, service.ErrCodeSendTooMany):
		ctx.JSON(http.StatusOK, Result{Code: 4, Msg: "验证码发送太频繁, 请稍后再试"})
	default:
		ctx.JSON(http.StatusOK, Result{Msg: "系统错误"})
	}
}

func (h *UserHandler) SMSLogin(ctx *gin.Context) {
	type SMSLoginReq struct {
		Phone string `json:"phone"`
		Code  string `json:"code"`
	}
	var req SMSLoginReq
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusOK, Result{Code: 5, Msg: "系统错误"})
		return
	}
	if req.Code == "" {
		ctx.JSON(http.StatusOK, Result{Code: 4, Msg: "验证码为空，请输入验证码"})
		return
	}
	ok, err := h.codeSvc.Verify(ctx, biz, req.Phone, req.Code)
	if errors.Is(err, service.ErrCodeVerifyTooMany) {
		// 可能有人搞你
		ctx.JSON(http.StatusOK, Result{Code: 6, Msg: "验证太频繁，请稍后再试"})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusOK, Result{Code: 5, Msg: "系统错误"})
		return
	}
	if !ok {
		ctx.JSON(http.StatusOK, Result{Code: 4, Msg: "验证码有误"})
		return
	}
	user, err := h.svc.FindOrCreate(ctx, req.Phone)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{Code: 5, Msg: "系统错误"})
		return
	}
	if err = h.userHdl.SetLoginToken(ctx, user.Id); err != nil {
		ctx.JSON(http.StatusOK, Result{Code: 5, Msg: "系统错误"})
		return
	}
	ctx.JSON(http.StatusOK, Result{Msg: "登陆成功"})
	return
}

func (h *UserHandler) jwtMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 1.从请求中获取 token
		tokenStr := extractTokenFromHeader(ctx)
		if tokenStr == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized,
				Result{Code: 4, Msg: "未提供访问令牌"})
			return
		}
		// 2.验证 JWT
		claims, err := h.userHdl.ParseToken(ctx, tokenStr)
		if err != nil {
			handleJWTError(ctx, err)
			ctx.Abort()
			return
		}
		// 3.检查 token 是否已过期
		if claims.ExpiresAt != nil {
			// 转换为 Unix 时间戳(秒级)
			expiresAt := claims.ExpiresAt.Unix()
			currentTime := time.Now().Unix()
			if currentTime > expiresAt {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, Result{Code: 4, Msg: "令牌已过期"})
				return
			} else {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, Result{Code: 4, Msg: "令牌格式异常"})
				return
			}
		}
		// 4.将 claims 存入上下文
		ctx.Set("uc", claims)
		ctx.Next()
	}
}

// 从 Header 中提取 token
func extractTokenFromHeader(ctx *gin.Context) string {
	authHeader := ctx.GetHeader("Authorization")
	if len(authHeader) < 8 || !strings.HasPrefix(authHeader, "Bearer ") {
		return ""
	}
	return authHeader[7:]
}

// 处理 JWT 验证错误
func handleJWTError(ctx *gin.Context, err error) {
	var result Result
	switch {
	case errors.Is(err, ijwt.ErrInvalidToken):
		result = Result{Code: 4, Msg: "无效的令牌"}
	case errors.Is(err, ijwt.ErrTokenExpired):
		result = Result{Code: 5, Msg: "令牌已过期"}
	default:
		result = Result{Code: 6, Msg: "令牌验证失败"}
	}
	ctx.AbortWithStatusJSON(http.StatusUnauthorized, result)
}
