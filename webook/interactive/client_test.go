package main

import (
	"context"
	"github.com/stretchr/testify/require"
	intrv1 "github.com/zmsocc/practice/webook/api/proto/gen/intr/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"testing"
)

func TestGRPCClient(t *testing.T) {
	cc, err := grpc.NewClient("localhost:8090",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	client := intrv1.NewInteractiveServiceClient(cc)
	res, err := client.Get(context.Background(), &intrv1.GetRequest{
		Biz:   "test",
		BizId: 2,
		Uid:   345,
	})
	require.NoError(t, err)
	t.Log(res.Intr)
}
