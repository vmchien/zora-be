package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

type ZnsReq struct {
	Hello string
}

type ZnsRepo interface {
	Save(context.Context, *ZnsReq) (*ZnsReq, error)
	Update(context.Context, *ZnsReq) (*ZnsReq, error)
	FindByID(context.Context, int64) (*ZnsReq, error)
	ListByHello(context.Context, string) ([]*ZnsReq, error)
	ListAll(context.Context) ([]*ZnsReq, error)
}

type ZnsUseCase struct {
	repo ZnsRepo
}

func NewZnsUseCase(repo ZnsRepo) *ZnsUseCase {
	return &ZnsUseCase{repo: repo}
}

func (uc *ZnsUseCase) Send(ctx context.Context, g *ZnsReq) (*ZnsReq, error) {
	log.Infof("Send: %v", g.Hello)
	return uc.repo.Save(ctx, g)
}
