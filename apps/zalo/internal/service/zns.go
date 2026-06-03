package service

import (
	"context"

	v1 "vn.vato.zora.be.api/api/zalo/v1"
	"vn.vato.zora.be.api/apps/zalo/internal/biz"
)

type ZnsService struct {
	v1.UnimplementedZnsServiceServer
	uc *biz.ZnsUseCase
}

func NewZnsService(uc *biz.ZnsUseCase) *ZnsService {
	return &ZnsService{uc: uc}
}

func (s *ZnsService) Send(ctx context.Context, in *v1.ZnsRequest) (*v1.ZnsReply, error) {
	g, err := s.uc.Send(ctx, &biz.ZnsReq{Hello: in.Name})
	if err != nil {
		return nil, err
	}
	return &v1.ZnsReply{Message: "Hello " + g.Hello}, nil
}
