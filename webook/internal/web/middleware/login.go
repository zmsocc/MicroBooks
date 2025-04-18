package middleware

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

type LoginMiddlewareBuild struct {
	paths []string
}

func NewLoginMiddlewareBuild() *LoginMiddlewareBuild {
	return &LoginMiddlewareBuild{}
}

func (l *LoginMiddlewareBuild) IgnorePaths(path string) *LoginMiddlewareBuild {
	l.paths = append(l.paths, path)
	return l
}

func (l *LoginMiddlewareBuild) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		for _, path := range l.paths {
			if ctx.Request.URL.Path == path {
				return
			}
		}
		sess := sessions.Default(ctx)
		if sess.Get("userId") == nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	}
}

func (l *LoginMiddlewareBuild) CheckLogin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if ctx.Request.URL.Path == "/users/signup" || ctx.Request.URL.Path == "/users/login" {
			return
		}

	}
}
