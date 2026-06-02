package data

import (
	"context"

	"vn.vato.zora.be.api/apps/zalo/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type ZnsRepo struct {
	data *Data
	log  *log.Helper
}

func NewZnsRepo(data *Data, logger log.Logger) biz.ZnsRepo {
	return &ZnsRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *ZnsRepo) Save(ctx context.Context, req *biz.ZnsReq) (*biz.ZnsReq, error) {
	// TODO implement me
	panic("implement me")
}

func (r *ZnsRepo) Update(ctx context.Context, req *biz.ZnsReq) (*biz.ZnsReq, error) {
	// TODO implement me
	panic("implement me")
}

func (r *ZnsRepo) FindByID(ctx context.Context, i int64) (*biz.ZnsReq, error) {
	// TODO implement me
	panic("implement me")
}

func (r *ZnsRepo) ListByHello(ctx context.Context, s string) ([]*biz.ZnsReq, error) {
	// TODO implement me
	panic("implement me")
}

func (r *ZnsRepo) ListAll(ctx context.Context) ([]*biz.ZnsReq, error) {
	// TODO implement me
	panic("implement me")
}
