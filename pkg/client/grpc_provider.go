package client

import (
	"context"
	"crypto/tls"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/middleware/metadata"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	ggrpc "google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/keepalive"
)

type GrpcClientOption struct {
	Timeout         time.Duration
	Mod             string
	RetryConnection bool
	MaxRetries      int32
	RetryInterval   time.Duration
}

const (
	defaultTimeout       = 30 * time.Second
	defaultRetryInterval = 2 * time.Second
	defaultMaxRetries    = 3

	// Conservative keepalive to avoid ENHANCE_YOUR_CALM: too_many_pings
	kaTime    = 5 * time.Minute
	kaTimeout = 20 * time.Second
)

func normalize(opt *GrpcClientOption) *GrpcClientOption {
	if opt == nil {
		opt = &GrpcClientOption{
			Timeout:         30 * time.Second,
			Mod:             "prod",
			RetryConnection: true,
			MaxRetries:      defaultMaxRetries,
			RetryInterval:   defaultRetryInterval,
		}
	}

	if opt.Timeout <= 0 {
		opt.Timeout = defaultTimeout
	}
	if opt.RetryConnection {
		if opt.RetryInterval <= 0 {
			opt.RetryInterval = defaultRetryInterval
		}
		if opt.MaxRetries <= 0 {
			opt.MaxRetries = defaultMaxRetries
		}
	}
	return opt
}

func dialOptions(endpoint string, opt *GrpcClientOption) []grpc.ClientOption {
	opts := []grpc.ClientOption{
		grpc.WithEndpoint(endpoint),
		grpc.WithTimeout(opt.Timeout), // per-RPC timeout in Kratos transport layer
		grpc.WithUnaryInterceptor(UnaryClientErrorInterceptor),
		grpc.WithMiddleware(
			recovery.Recovery(),
			tracing.Client(),
			metadata.Client(),
		),
	}

	if opt.Mod == "local" && strings.Contains(endpoint, ":443") {
		opts = append(opts, grpc.WithTLSConfig(&tls.Config{InsecureSkipVerify: true}))
	}

	if opt.RetryConnection {
		keepaliveParams := keepalive.ClientParameters{
			Time:                kaTime,
			Timeout:             kaTimeout,
			PermitWithoutStream: false, // avoid pings when idle → prevents too_many_pings
		}
		opts = append(opts, grpc.WithOptions(
			ggrpc.WithKeepaliveParams(keepaliveParams),
			ggrpc.WithConnectParams(ggrpc.ConnectParams{
				Backoff:           backoff.Config{MaxDelay: opt.RetryInterval},
				MinConnectTimeout: opt.Timeout,
			}),
		))
	}
	return opts
}

func ensureDNSResolverTarget(endpoint string) string {
	if strings.HasPrefix(endpoint, "dns:///") {
		return endpoint
	}
	if strings.Contains(endpoint, "://") {
		return endpoint
	}
	return "dns:///" + endpoint
}

// // CreateConnection creates a gRPC connection without using deprecated dial options.
// // It attempts to establish a READY connection within opt.Timeout per attempt and retries when configured.
// func CreateConnectionV0(endpoint string, opt *GrpcClientOption) *ggrpc.ClientConn {
// 	opt = normalize(opt)
// 	opts := dialOptions(endpoint, opt)
//
// 	attempts := int32(1)
// 	if opt.RetryConnection && opt.MaxRetries > 1 {
// 		attempts = opt.MaxRetries
// 	}
//
// 	var (
// 		conn *ggrpc.ClientConn
// 		err  error
// 	)
//
// 	for i := int32(1); i <= attempts; i++ {
// 		// Non-blocking dial (no deprecated options)
// 		conn, err = grpc.DialInsecure(context.Background(), opts...)
// 		if err != nil {
// 			if i < attempts {
// 				time.Sleep(opt.RetryInterval)
// 				continue
// 			}
// 			break
// 		}
//
// 		// Actively try to connect and wait until READY or timeout per attempt
// 		attemptCtx, cancel := context.WithTimeout(context.Background(), opt.Timeout)
// 		conn.Connect()
//
// 		ready := false
// 		for {
// 			state := conn.GetState()
// 			if state == connectivity.Ready {
// 				ready = true
// 				break
// 			}
// 			if !conn.WaitForStateChange(attemptCtx, state) {
// 				// Timeout or context canceled
// 				break
// 			}
// 		}
// 		cancel()
//
// 		if ready {
// 			return conn
// 		}
//
// 		_ = conn.Close()
// 		conn = nil
// 		if i < attempts {
// 			time.Sleep(opt.RetryInterval)
// 		}
// 	}
//
// 	panic(fmt.Errorf("failed to establish READY connection to %s after %d attempt(s): %v", endpoint, attempts, err))
// }

func CreateConnection(endpoint string, opt *GrpcClientOption) *ggrpc.ClientConn {
	opt = normalize(opt)
	target := ensureDNSResolverTarget(endpoint)

	opts := []grpc.ClientOption{
		grpc.WithEndpoint(target),
		grpc.WithTimeout(opt.Timeout),
		// grpc.WithUnaryInterceptor(UnaryClientErrorInterceptor), // TODO: check and remove if not needed after error handling is implemented in middleware
		grpc.WithMiddleware(
			recovery.Recovery(),
			tracing.Client(),
			metadata.Client(),
		),

		// Implement round_robin with headless service
		grpc.WithOptions(
			ggrpc.WithDefaultServiceConfig(`{"loadBalancingConfig":[{"round_robin":{}}]}`),
			ggrpc.WithConnectParams(ggrpc.ConnectParams{
				Backoff: backoff.Config{
					MaxDelay: opt.RetryInterval,
				},
				MinConnectTimeout: opt.Timeout,
			}),
			ggrpc.WithKeepaliveParams(keepalive.ClientParameters{
				Time:                kaTime,
				Timeout:             kaTimeout,
				PermitWithoutStream: false,
			}),
		),
	}

	if opt.Mod == "local" && strings.Contains(endpoint, ":443") {
		opts = append(opts, grpc.WithTLSConfig(&tls.Config{InsecureSkipVerify: true}))
	}

	conn, err := grpc.DialInsecure(context.Background(), opts...)
	if err != nil {
		panic(err)
	}

	return conn
}
