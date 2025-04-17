package web

import (
	"errors"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/zmsocc/practice/webook/internal/domain"
	"github.com/zmsocc/practice/webook/internal/service"
	"net/http"
	"regexp"
	"strings"
	"time"
)

const (
	emailRegexpPattern    = "^[a-zA-Z0-9_-]+@[a-zA-Z0-9_-]+(\\.[a-zA-Z0-9_-]+)+$"
	passwordRegexpPattern = "^\\d{1,9}$" // 只能由数字组成
)

// ^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$

type UserHandler struct {
	svc *service.UserService
	//emailRegexp    *regexp.Regexp
	//passwordRegexp *regexp.Regexp
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{
		svc: svc,
		//emailRegexp:    regexp.MustCompile(emailRegexpPattern, regexp.None),
		//passwordRegexp: regexp.MustCompile(passwordRegexpPattern, regexp.None),
	}
}

func InitWebServer() *gin.Engine {
	server := gin.Default()
	server.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowHeaders:     []string{"Content-Type"},
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				return true
			}
			return strings.Contains(origin, "your_company.com")
		},
		MaxAge: 12 * time.Hour,
	}))
	return server
}

func (h *UserHandler) RegisterRoutes(server *gin.Engine) {
	ug := server.Group("/users")
	ug.POST("/signup", h.SignUp)
	//ug.POST("/login", u.Login)
	//ug.POST("/edit", u.Edit)
	//ug.GET("/profile", u.Profile)
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
	isEmail, err := regexp.Match(emailRegexpPattern, []byte(req.Email))
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
	isPassword, err := regexp.Match(passwordRegexpPattern, []byte(req.Password))
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

//func (u *UserHandler) Login(ctx *gin.Context) {
//
//}
//
//func (u *UserHandler) Edit(ctx *gin.Context) {
//
//}
//
//func (u *UserHandler) Profile(ctx *gin.Context) {
//
//}
