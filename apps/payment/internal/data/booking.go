package data

import (
	"context"

	"vn.vato.zora.be.api/apps/payment/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type bookingRepo struct {
	data *Data
	log  *log.Helper
}

// NewBookingRepo .
func NewBookingRepo(data *Data, logger log.Logger) biz.BookingRepo {
	return &bookingRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *bookingRepo) Save(ctx context.Context, g *biz.Booking) (*biz.Booking, error) {
	return g, nil
}

func (r *bookingRepo) Update(ctx context.Context, g *biz.Booking) (*biz.Booking, error) {
	return g, nil
}

func (r *bookingRepo) FindByID(context.Context, int64) (*biz.Booking, error) {
	return nil, nil
}

func (r *bookingRepo) ListByHello(context.Context, string) ([]*biz.Booking, error) {
	return nil, nil
}

func (r *bookingRepo) ListAll(context.Context) ([]*biz.Booking, error) {
	return nil, nil
}
