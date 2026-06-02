package data

import (
	"context"

	zaloV1 "vn.vato.zora.be.api/api/zalo/v1"
	"vn.vato.zora.be.api/apps/gateway/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type PaymentClient struct {
	client zaloV1.ZnsServiceClient
	log    *log.Helper
}

// NewPaymentRepo new a PaymentRepo implementation.
func NewPaymentRepo(client zaloV1.ZnsServiceClient, logger log.Logger) biz.PaymentRepo {
	return &PaymentClient{
		client: client,
		log:    log.NewHelper(logger),
	}
}

func (r *PaymentClient) BookingTicket(ctx context.Context, name string) (string, error) {
	resp, err := r.client.Send(ctx, &zaloV1.ZnsRequest{Name: name})
	if err != nil {
		return "", err
	}
	return resp.Message, nil
}
