package data

import (
	"context"

	paymentv1 "vn.vato.zora.be.api/api/payment/v1"
	walletv1 "vn.vato.zora.be.api/api/wallet/v1"
	"vn.vato.zora.be.api/apps/gateway/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/google/wire"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewGreeterRepo, NewPaymentClient, NewPaymentRepo, NewWalletRepo, NewWalletClient)

// Data .
type Data struct {
	// TODO wrapped database client
}

// NewData .
func NewData(c *conf.Data) (*Data, func(), error) {
	cleanup := func() {
		log.Info("closing the data resources")
	}
	return &Data{}, cleanup, nil
}

// NewPaymentClient new a gRPC client for payment service.
func NewPaymentClient(c *conf.Data, logger log.Logger) (paymentv1.BookingServiceClient, error) {
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint(c.Payment.Endpoint),
		grpc.WithTimeout(c.Payment.Timeout.AsDuration()),
	)
	if err != nil {
		return nil, err
	}
	return paymentv1.NewBookingServiceClient(conn), nil
}

func NewWalletClient(c *conf.Data, logger log.Logger) (walletv1.WalletServiceClient, error) {
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint(c.Wallet.Endpoint),
		grpc.WithTimeout(c.Wallet.Timeout.AsDuration()),
	)
	if err != nil {
		return nil, err
	}
	return walletv1.NewWalletServiceClient(conn), nil
}
