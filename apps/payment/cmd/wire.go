//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"vn.vato.zora.be.api/apps/payment/internal/biz"
	"vn.vato.zora.be.api/apps/payment/internal/conf"
	"vn.vato.zora.be.api/apps/payment/internal/data"
	"vn.vato.zora.be.api/apps/payment/internal/server"
	"vn.vato.zora.be.api/apps/payment/internal/service"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(*conf.Server, *conf.Data, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, newApp))
}
