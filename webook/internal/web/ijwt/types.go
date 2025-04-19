package ijwt

import "github.com/gin-gonic/gin"

type Handler interface {
	CheckSession(ctx *gin.Context, ssid string) error
	SetJWTToken(ctx *gin.Context, uid int64, ssid string) error
	SetLoginToken(ctx *gin.Context, uid int64) error
	SetRefreshToken(ctx *gin.Context, uid int64, ssid string) error
}
