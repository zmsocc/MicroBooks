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
)

const (
	emailRegexpPattern    = "^[a-zA-Z0-9_-]+@[a-zA-Z0-9_-]+(\\.[a-zA-Z0-9_-]+)+$"
	passwordRegexpPattern = "^(?=.*[A-Za-z])(?=.*\\d)(?=.*[$@$!%*#?&])[A-Za-z\\d$@$!%*#?&]{8,}$" // 只能由数字组成
)

// ^\d{1,9}$

type UserHandler struct {
	svc            *service.UserService
	emailRegexp    *regexp.Regexp
	passwordRegexp *regexp.Regexp
	userHdl        ijwt.Handler
}

func NewUserHandler(svc *service.UserService, userHdl ijwt.Handler) *UserHandler {
	return &UserHandler{
		svc:            svc,
		emailRegexp:    regexp.MustCompile(emailRegexpPattern, regexp.None),
		passwordRegexp: regexp.MustCompile(passwordRegexpPattern, regexp.None),
		userHdl:        userHdl,
	}
}

func (h *UserHandler) RegisterRoutes(server *gin.Engine) {
	ug := server.Group("/users")
	ug.POST("/signup", h.SignUp)
	ug.POST("/login", ginx.WrapBodyV1(h.LoginJWT))
	//ug.POST("/edit", u.Edit)
	ug.GET("/profile", h.Profile)
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

func (h *UserHandler) Profile(ctx *gin.Context) {
	ctx.String(http.StatusOK, "这是你的 profile")
}
