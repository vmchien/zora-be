package service

import (
	"context"

	v1 "vn.vato.zora.be.api/api/payment/v1"
	"vn.vato.zora.be.api/apps/payment/internal/biz"
)

// PaymentService is a payment service.
type PaymentService struct {
	v1.UnimplementedBookingServiceServer

	uc *biz.BookingUsecase
}

// NewPaymentService new a payment service.
func NewPaymentService(uc *biz.BookingUsecase) *PaymentService {
	return &PaymentService{uc: uc}
}

// BookingTicket implements BookingServiceServer.
func (s *PaymentService) BookingTicket(ctx context.Context, in *v1.HelloRequest) (*v1.HelloReply, error) {
	g, err := s.uc.CreateBooking(ctx, &biz.Booking{Hello: in.Name})
	if err != nil {
		return nil, err
	}
	return &v1.HelloReply{Message: "Hello " + g.Hello}, nil
}
