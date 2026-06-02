package data

import (
	"context"

	paymentv1 "vn.vato.zora.be.api/api/payment/v1"
	"vn.vato.zora.be.api/apps/gateway/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type PaymentClient struct {
	client paymentv1.BookingServiceClient
	log    *log.Helper
}

// NewPaymentRepo new a PaymentRepo implementation.
func NewPaymentRepo(client paymentv1.BookingServiceClient, logger log.Logger) biz.PaymentRepo {
	return &PaymentClient{
		client: client,
		log:    log.NewHelper(logger),
	}
}

func (r *PaymentClient) BookingTicket(ctx context.Context, name string) (string, error) {
	resp, err := r.client.BookingTicket(ctx, &paymentv1.HelloRequest{Name: name})
	if err != nil {
		return "", err
	}
	return resp.Message, nil
}
