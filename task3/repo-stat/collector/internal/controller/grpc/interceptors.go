package grpccontroller

import (
	"context"
	"log/slog"

	"google.golang.org/grpc"
)

func LoggingInterceptor(log *slog.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		resp, err := handler(ctx, req)

		attributes := []slog.Attr{
			slog.String("method", info.FullMethod),
		}

		if err != nil {
			finalAttrs := append(attributes, slog.String("error", err.Error()))

			log.LogAttrs(ctx, slog.LevelError, "grpc request failed", finalAttrs...)
		} else {
			log.LogAttrs(ctx, slog.LevelInfo, "grpc request success", attributes...)
		}

		return resp, err
	}
}
