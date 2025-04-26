package tencent

import (
	"context"
	"fmt"
	"github.com/ecodeclub/ekit"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
)

type Service struct {
	client   *sms.Client
	appId    *string
	signName *string
}

func NewService(client *sms.Client, appId string, signName string) *Service {
	return &Service{
		client:   client,
		appId:    ekit.ToPtr[string](appId),
		signName: ekit.ToPtr[string](signName),
	}
}

func (s *Service) Send(ctx context.Context, biz string, args []string, numbers ...string) error {
	req := sms.NewSendSmsRequest()
	req.SmsSdkAppId = s.appId
	req.SignName = s.signName
	req.TemplateId = ekit.ToPtr[string](biz)
	req.PhoneNumberSet = toStringPtrSlice(numbers)
	req.TemplateParamSet = toStringPtrSlice(args)
	//req.SetContext(ctx)
	resp, err := s.client.SendSms(req)
	if err != nil {
		return fmt.Errorf("腾讯短信服务发送失败 %w", err)
	}
	for _, status := range resp.Response.SendStatusSet {
		if status.Code == nil || *(status.Code) != "Ok" {
			return fmt.Errorf("发送失败，code: %s, 原因: %s", *status.Code, *status.Message)
		}
	}
	return nil
}

func toStringPtrSlice(src []string) []*string {
	result := make([]*string, len(src))
	for i := range src {
		val := src[i]    // 创建局部变量副本
		result[i] = &val // 指向该副本的地址
	}
	return result
}
