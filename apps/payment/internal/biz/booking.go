package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

type Booking struct {
	Hello string
}

// BookingRepo is a Greater repo.
type BookingRepo interface {
	Save(context.Context, *Booking) (*Booking, error)
	Update(context.Context, *Booking) (*Booking, error)
	FindByID(context.Context, int64) (*Booking, error)
	ListByHello(context.Context, string) ([]*Booking, error)
	ListAll(context.Context) ([]*Booking, error)
}

// BookingUsecase is a Booking usecase.
type BookingUsecase struct {
	repo BookingRepo
}

// NewBookingUsecase new a Booking usecase.
func NewBookingUsecase(repo BookingRepo) *BookingUsecase {
	return &BookingUsecase{repo: repo}
}

// CreateBooking creates a Booking, and returns the new Booking.
func (uc *BookingUsecase) CreateBooking(ctx context.Context, g *Booking) (*Booking, error) {
	log.Infof("CreateBooking: %v", g.Hello)
	return uc.repo.Save(ctx, g)
}
