package data

import (
	"context"

	walletv1 "vn.vato.zora.be.api/api/wallet/v1"
	zaloV1 "vn.vato.zora.be.api/api/zalo/v1"

	"vn.vato.zora.be.api/apps/gateway/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/google/wire"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewGreeterRepo, NewZaloClient, NewPaymentRepo, NewWalletRepo, NewWalletClient)

type Data struct {
	// TODO wrapped database client
}

func NewData(c *conf.Data) (*Data, func(), error) {
	cleanup := func() {
		log.Info("closing the data resources")
	}
	return &Data{}, cleanup, nil
}

func NewZaloClient(c *conf.Data, logger log.Logger) (zaloV1.ZnsServiceClient, error) {
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint(c.Payment.Endpoint),
		grpc.WithTimeout(c.Payment.Timeout.AsDuration()),
	)
	if err != nil {
		return nil, err
	}
	return zaloV1.NewZnsServiceClient(conn), nil
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
