package cache

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
)

var (
	ErrCodeSendTooMany   = errors.New("发送太频繁")
	ErrCodeVerifyTooMany = errors.New("验证太频繁")
	ErrUnknownForCode    = errors.New("未知错误")
)

//go:embed lua/set_code.lua
var luaSetCode string

//go:embed lua/verify_code.lua
var luaVerifyCode string

type CodeCache interface {
	Get(ctx context.Context, biz, phone, inputCode string) (bool, error)
	Set(ctx context.Context, biz, phone, code string) error
}

type RedisCodeCache struct {
	cmd redis.Cmdable
}

//func NewRedisCodeCache(client *redis.Client) *RedisCodeCache {
//	return &RedisCodeCache{
//		client: client,
//	}
//}

func NewCodeCache(cmd redis.Cmdable) CodeCache {
	return &RedisCodeCache{
		cmd: cmd,
	}
}

func (c *RedisCodeCache) Get(ctx context.Context, biz, phone, inputCode string) (bool, error) {
	res, err := c.cmd.Eval(ctx, luaVerifyCode, []string{c.key(biz, phone)}, inputCode).Int()
	if err != nil {
		return false, err
	}
	switch res {
	case 0:
		// 验证码输入对了
		return true, nil
	case -1:
		// 要注意，可能有人恶意攻击
		return false, ErrCodeVerifyTooMany
	case -2:
		return false, nil
	default:
		return false, ErrUnknownForCode
	}
}

func (c *RedisCodeCache) Set(ctx context.Context, biz, phone, code string) error {
	// 确保 code 是字符串类型
	// log.Printf("传递参数类型：code=%T", code)
	res, err := c.cmd.Eval(ctx, luaSetCode, []string{c.key(biz, phone)}, code).Int()
	if err != nil {
		return err
	}
	switch res {
	case 0:
		// 毫无问题
		return nil
	case -1:
		return ErrCodeSendTooMany
	default:
		return errors.New("系统错误")
	}
}

func (c *RedisCodeCache) key(biz string, phone string) string {
	return fmt.Sprintf("phone_code:%s:%s", biz, phone)
}
