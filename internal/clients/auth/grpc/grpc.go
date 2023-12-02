package grpc

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	gspv1 "github.com/commedesvlados/grpc-service-protos/gen/go/grpc_service"
	grpclog "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	grpcretry "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthClient struct {
	api gspv1.AuthClient
}

func New(
	ctx context.Context,
	log *slog.Logger,
	addr string,
	timeout time.Duration,
	retriesCount int,
) (*AuthClient, error) {
	const fn = "clients.auth.grpc.New"

	retryOpts := []grpcretry.CallOption{
		grpcretry.WithCodes(codes.NotFound, codes.Aborted, codes.DeadlineExceeded),
		grpcretry.WithMax(uint(retriesCount)),
		grpcretry.WithPerRetryTimeout(timeout),
	}

	logOpts := []grpclog.Option{
		grpclog.WithLogOnEvents(grpclog.PayloadReceived, grpclog.PayloadSent),
	}

	cc, err := grpc.DialContext(ctx, addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()), // TODO secure
		grpc.WithChainUnaryInterceptor(
			grpclog.UnaryClientInterceptor(InterceptorLogger(log), logOpts...),
			grpcretry.UnaryClientInterceptor(retryOpts...),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return &AuthClient{
		api: gspv1.NewAuthClient(cc),
	}, nil
}

func InterceptorLogger(l *slog.Logger) grpclog.Logger {
	return grpclog.LoggerFunc(func(ctx context.Context, level grpclog.Level, msg string, fields ...any) {
		l.Log(ctx, slog.Level(level), msg, fields...)
	})
}

func (a *AuthClient) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	const fn = "clients.auth.grpc.IsAdmin"

	resp, err := a.api.IsAdmin(ctx, &gspv1.IsAdminRequest{UserId: userID})
	if err != nil {
		return false, fmt.Errorf("%s: %w", fn, err)
	}

	return resp.GetIsAdmin(), nil
}
