package ijwt

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/zmsocc/practice/webook/pkg/ginx"
	"net/http"
	"os"
	"time"
)

// AtKey        = []byte("95osj3fUD7fo0mlYdDbncXz4VD2igvf0")
//RtKey        = []byte("0Pf2r0wZBpXVXlQNdpwCXN4ncnlnZSc3")

var (
	atPrivateKey *ecdsa.PrivateKey
)

type RedisJWTHandler struct {
	cmd redis.Cmdable
}

func NewRedisJWTHandler(cmd redis.Cmdable) Handler {
	return &RedisJWTHandler{
		cmd: cmd,
	}
}

func init() {
	// 加载 Access Token 密钥
	privkey, _ := os.ReadFile("ec512-private.pem")
	block, _ := pem.Decode(privkey)
	key, _ := x509.ParseECPrivateKey(block.Bytes)
	atPrivateKey = key
}

func (h *RedisJWTHandler) CheckSession(ctx *gin.Context, ssid string) error {
	val, err := h.cmd.Exists(ctx, fmt.Sprintf("users:ssid:%s", ssid)).Result()
	switch {
	case errors.Is(err, redis.Nil):
		return nil
	case err == nil:
		if val == 0 {
			return nil
		}
		return errors.New("session 已经失效了")
	default:
		return err
	}
}

func (h *RedisJWTHandler) SetJWTToken(ctx *gin.Context, uid int64, ssid string) error {
	uc := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			// 过期时间设置为 1 分钟, 测试
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute)),
		},
		Uid:       uid,
		Ssid:      ssid,
		UserAgent: ctx.Request.UserAgent(),
	}
	// 使用 ECDSA 密钥签名
	token := jwt.NewWithClaims(jwt.SigningMethodES512, uc)
	// 使用正确格式的私钥
	tokenStr, err := token.SignedString(atPrivateKey)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ginx.Result{
			Code: 5,
			Msg:  "令牌生成失败",
			Data: err.Error(),
		})
		return err
	}
	ctx.Header("X-Jwt-Token", tokenStr)
	return nil
}

func (h *RedisJWTHandler) SetLoginToken(ctx *gin.Context, uid int64) error {
	ssid := uuid.New().String()
	err := h.SetJWTToken(ctx, uid, ssid)
	if err != nil {
		return err
	}
	err = h.SetRefreshToken(ctx, uid, ssid)
	return err
}

func (h *RedisJWTHandler) SetRefreshToken(ctx *gin.Context, uid int64, ssid string) error {
	uc := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute)),
		},
		Uid:       uid,
		Ssid:      ssid,
		UserAgent: ctx.Request.UserAgent(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodES512, uc)
	tokenStr, err := token.SignedString(atPrivateKey)
	if err != nil {
		ctx.Header("x-refresh-token", tokenStr)
		return nil
	}
	return err
}

func (h *RedisJWTHandler) ParseToken(ctx *gin.Context, tokenStr string) (UserClaims, error) {
	// 解析 JWT 令牌
	var uc UserClaims
	token, err := jwt.ParseWithClaims(tokenStr, &uc, h.keyFunc)
	if err != nil {
		return UserClaims{}, err
	}
	if !token.Valid {
		return UserClaims{}, ErrInvalidToken
	}
	// Redis 会话校验
	if err = h.CheckSession(ctx, uc.Ssid); err != nil {
		return UserClaims{}, fmt.Errorf("%w: session invalid", ErrInvalidToken)
	}
	return uc, nil
}

// 密钥处理函数
func (h *RedisJWTHandler) keyFunc(token *jwt.Token) (interface{}, error) {
	switch token.Method {
	case jwt.SigningMethodES512:
		return atPrivateKey, nil
	default:
		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])

	}
}

type UserClaims struct {
	jwt.RegisteredClaims
	Uid       int64
	Ssid      string
	UserAgent string
}
