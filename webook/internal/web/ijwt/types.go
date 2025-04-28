package ijwt

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrTokenExpired = errors.New("token is expired")
)

type Handler interface {
	CheckSession(ctx *gin.Context, ssid string) error
	ClearToken(ctx *gin.Context) error
	ExtractToken(ctx *gin.Context) string
	SetJWTToken(ctx *gin.Context, uid int64, ssid string) error
	SetLoginToken(ctx *gin.Context, uid int64) error
	SetRefreshToken(ctx *gin.Context, uid int64, ssid string) error
	ParseToken(ctx *gin.Context, tokenStr string) (UserClaims, error)
}

type RefreshClaims struct {
	Uid  int64
	Ssid string
	jwt.RegisteredClaims
}

type UserClaims struct {
	jwt.RegisteredClaims
	Uid       int64
	Ssid      string
	UserAgent string
}
