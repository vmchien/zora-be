package data

import (
	"context"

	zaloV1 "vn.vato.zora.be.api/api/zalo/v1"
	"vn.vato.zora.be.api/apps/gateway/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(NewData, NewGreeterRepo, NewZaloClient, NewPaymentRepo)

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
