package ijwt

import (
	"errors"
	"github.com/gin-gonic/gin"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrTokenExpired = errors.New("token is expired")
)

type Handler interface {
	CheckSession(ctx *gin.Context, ssid string) error
	SetJWTToken(ctx *gin.Context, uid int64, ssid string) error
	SetLoginToken(ctx *gin.Context, uid int64) error
	SetRefreshToken(ctx *gin.Context, uid int64, ssid string) error
	ParseToken(ctx *gin.Context, tokenStr string) (UserClaims, error)
}
