package client

import (
	"context"

	"google.golang.org/grpc"
	errWrap "vn.vato.zora.be.api/pkg/errors"
)

func GrpcErrorHandler(_ context.Context, err error) error {
	if e := errWrap.FromError(err); e != nil {
		return errWrap.ToStatusError(e)
	}
	return err
}
func UnaryServerErrorInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	resp, err = handler(ctx, req)
	if err == nil {
		return resp, nil
	}
	if e := errWrap.FromError(err); e != nil {
		return nil, errWrap.ToStatusError(e)
	}
	return nil, err
}

func UnaryClientErrorInterceptor(ctx context.Context, method string, req, reply any,
	cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	err := invoker(ctx, method, req, reply, cc, opts...)
	if err == nil {
		return nil
	}
	if e := errWrap.FromStatusError(err); e != nil {
		return e
	}
	return err
}
