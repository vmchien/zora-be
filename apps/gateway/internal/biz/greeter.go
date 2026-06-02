package biz

import (
	"context"

	v1 "vn.vato.zora.be.api/api/gateway/v1"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

var (
	// ErrUserNotFound is user not found.
	ErrUserNotFound = errors.NotFound(v1.ErrorReason_USER_NOT_FOUND.String(), "user not found")
)

// Greeter is a Greeter model.
type Greeter struct {
	Hello string
}

// GreeterRepo is a Greater repo.
type GreeterRepo interface {
	Save(context.Context, *Greeter) (*Greeter, error)
	Update(context.Context, *Greeter) (*Greeter, error)
	FindByID(context.Context, int64) (*Greeter, error)
	ListByHello(context.Context, string) ([]*Greeter, error)
	ListAll(context.Context) ([]*Greeter, error)
}

// PaymentRepo is a Payment client interface.
type PaymentRepo interface {
	BookingTicket(ctx context.Context, name string) (string, error)
}

// GreeterUsecase is a Greeter usecase.
type GreeterUsecase struct {
	repo GreeterRepo
	pay  PaymentRepo
}

// NewGreeterUsecase new a Greeter usecase.
func NewGreeterUsecase(repo GreeterRepo, pay PaymentRepo) *GreeterUsecase {
	return &GreeterUsecase{repo: repo, pay: pay}
}

// CreateGreeter creates a Greeter, and returns the new Greeter.
func (uc *GreeterUsecase) CreateGreeter(ctx context.Context, g *Greeter) (*Greeter, error) {
	log.Infof("CreateGreeter: %v", g.Hello)
	payMsg, err := uc.pay.BookingTicket(ctx, g.Hello)
	if err != nil {
		log.Errorf("failed to call payment service: %v", err)
		return nil, err
	}
	return uc.repo.Save(ctx, &Greeter{Hello: g.Hello + " (payment call result: " + payMsg + ")"})
}
