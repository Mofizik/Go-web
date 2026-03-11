package grpcmw

import (
	"context"
	"log/slog"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func UnaryServerLoggingInterceptor(log *slog.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		start := time.Now()
		resp, err = handler(ctx, req)
	
		code := status.Code(err)
		log.Info(
			"method", info.FullMethod,
			"req", req,
			"resp", resp,
			"duration", time.Since(start),
			"code", code.String(),
			"err", err,
		)
		return resp, err
	}
}