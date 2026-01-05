package ioc

import "github.com/zmsocc/practice/webook/pkg/logger"

func InitLogger() logger.Logger {
	return logger.NewNopLogger()
}
