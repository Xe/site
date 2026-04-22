package thoth

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/timeout"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	healthv1 "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/metadata"
	adminv1 "xeiaso.net/v4/gen/techaro/thoth/auth/admin/v1"
	authv1 "xeiaso.net/v4/gen/techaro/thoth/auth/v1"
)

type Client struct {
	conn *grpc.ClientConn

	Health     healthv1.HealthClient
	AuthJWT    authv1.JWTServiceClient
	AdminUsers adminv1.UsersServiceClient
}

func New(ctx context.Context, thothURL, apiToken string, noTLS bool) (*Client, error) {
	clMetrics := grpcprom.NewClientMetrics(
		grpcprom.WithClientHandlingTimeHistogram(
			grpcprom.WithHistogramBuckets([]float64{0.001, 0.01, 0.1, 0.3, 0.6, 1, 3, 6, 9, 20, 30, 60, 90, 120}),
		),
	)
	prometheus.DefaultRegisterer.Register(clMetrics)

	var transportCreds credentials.TransportCredentials

	switch noTLS {
	case true:
		transportCreds = insecure.NewCredentials()
	case false:
		transportCreds = credentials.NewTLS(&tls.Config{})
	}

	conn, err := grpc.NewClient(
		thothURL,
		grpc.WithTransportCredentials(transportCreds),
		grpc.WithChainUnaryInterceptor(
			timeout.UnaryClientInterceptor(5*time.Minute),
			clMetrics.UnaryClientInterceptor(),
			authUnaryClientInterceptor(apiToken),
		),
		grpc.WithChainStreamInterceptor(
			clMetrics.StreamClientInterceptor(),
			authStreamClientInterceptor(apiToken),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("can't dial thoth at %s: %w", thothURL, err)
	}

	hc := healthv1.NewHealthClient(conn)

	resp, err := hc.Check(ctx, &healthv1.HealthCheckRequest{})
	if err != nil {
		return nil, fmt.Errorf("can't verify thoth health at %s: %w", thothURL, err)
	}

	if resp.Status != healthv1.HealthCheckResponse_SERVING {
		return nil, fmt.Errorf("thoth is not healthy, wanted %s but got %s", healthv1.HealthCheckResponse_SERVING, resp.Status)
	}

	return &Client{
		conn:       conn,
		Health:     hc,
		AuthJWT:    authv1.NewJWTServiceClient(conn),
		AdminUsers: adminv1.NewUsersServiceClient(conn),
	}, nil
}

func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

func authUnaryClientInterceptor(token string) grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req any,
		reply any,
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		md := metadata.Pairs("authorization", "Bearer "+token)
		ctx = metadata.NewOutgoingContext(ctx, md)
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func authStreamClientInterceptor(token string) grpc.StreamClientInterceptor {
	return func(
		ctx context.Context,
		desc *grpc.StreamDesc,
		cc *grpc.ClientConn,
		method string,
		streamer grpc.Streamer,
		opts ...grpc.CallOption,
	) (grpc.ClientStream, error) {
		md := metadata.Pairs("authorization", "Bearer "+token)
		ctx = metadata.NewOutgoingContext(ctx, md)
		return streamer(ctx, desc, cc, method, opts...)
	}
}
