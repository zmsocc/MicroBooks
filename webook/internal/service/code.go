package service

import (
	"context"
	"fmt"
	"github.com/zmsocc/practice/webook/internal/repository"
	"github.com/zmsocc/practice/webook/internal/service/sms"
	"go.uber.org/atomic"
	"log"
	"math/rand"
)

var (
	ErrCodeSendTooMany   = repository.ErrCodeSendTooMany
	codeTplId            = atomic.String{}
	ErrCodeVerifyTooMany = repository.ErrCodeVerifyTooMany
)

type CodeService interface {
	Send(ctx context.Context, biz, phone string) error
	Verify(ctx context.Context, biz, phone, code string) (bool, error)
}

type codeService struct {
	repo   repository.CodeRepository
	smsSvc sms.Service
}

func NewCodeService(repo repository.CodeRepository, smsSvc sms.Service) CodeService {
	codeTplId.Store("1877556")
	return &codeService{
		repo:   repo,
		smsSvc: smsSvc,
	}
}

// Send -biz 区别业务场景
func (svc *codeService) Send(ctx context.Context, biz string, phone string) error {
	code := svc.generateCode()
	if err := svc.repo.Store(ctx, biz, phone, code); err != nil {
		log.Printf("验证码存储失败：biz=%s, phone=%s, err=%v", biz, phone, err)
	}
	err := svc.smsSvc.Send(ctx, codeTplId.Load(), []string{code}, phone)
	if err != nil {
		return fmt.Errorf("发送短信出现异常 %w", err)
	}
	return err
}

func (svc *codeService) Verify(ctx context.Context, biz, phone, code string) (bool, error) {
	return svc.repo.Verify(ctx, biz, phone, code)
}

func (svc *codeService) generateCode() string {
	randNum := rand.Intn(1000000)
	// 不够六位的，前面加上0，补够六位
	return fmt.Sprintf("%06d", randNum)
}
