package grpccontroller

import (
	"errors"

	"github.com/IliaSotnikov2005/golang-course/task3/repo-stat/collector/internal/domain"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func mapErrorToGRPC(err error) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, domain.ErrNotFound):
		return status.Error(codes.NotFound, err.Error())

	case errors.Is(err, domain.ErrMovedPermanently):
		return status.Error(codes.NotFound, err.Error())

	case errors.Is(err, domain.ErrForbidden):
		return status.Error(codes.PermissionDenied, err.Error())

	case errors.Is(err, domain.ErrUnauthorized):
		return status.Error(codes.Unauthenticated, err.Error())

	case errors.Is(err, domain.ErrRateLimit):
		return status.Error(codes.ResourceExhausted, err.Error())

	case errors.Is(err, domain.ErrInvalidInput):
		return status.Error(codes.InvalidArgument, err.Error())

	case errors.Is(err, domain.ErrTimeout):
		return status.Error(codes.DeadlineExceeded, err.Error())

	case errors.Is(err, domain.ErrInternal):
		return status.Error(codes.Internal, err.Error())

	default:
		return status.Error(codes.Internal, "internal server error")
	}
}
