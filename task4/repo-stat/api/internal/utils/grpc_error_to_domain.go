package utils

import (
	"fmt"

	"github.com/IliaSotnikov2005/golang-course/task4/repo-stat/api/internal/domain"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func MapGRPCErrorToDomain(err error) error {
	if err == nil {
		return nil
	}

	st, ok := status.FromError(err)
	if !ok {
		return fmt.Errorf("%w: %v", domain.ErrInternal, err)
	}

	switch st.Code() {
	case codes.AlreadyExists:
		return fmt.Errorf("%w: %s", domain.ErrSubscriptionAlreadyExists, st.Message())
	case codes.NotFound:
		return fmt.Errorf("%w: %s", domain.ErrNotFound, st.Message())

	case codes.PermissionDenied:
		return fmt.Errorf("%w: %s", domain.ErrForbidden, st.Message())

	case codes.Unauthenticated:
		return fmt.Errorf("%w: %s", domain.ErrUnauthorized, st.Message())

	case codes.ResourceExhausted:
		return fmt.Errorf("%w: %s", domain.ErrRateLimit, st.Message())

	case codes.InvalidArgument:
		return fmt.Errorf("%w: %s", domain.ErrInvalidInput, st.Message())

	case codes.DeadlineExceeded:
		return fmt.Errorf("%w: request to collector service timed out", domain.ErrTimeout)

	case codes.Unavailable:
		return fmt.Errorf("%w: collector service is unavailable", domain.ErrInternal)

	case codes.Internal:
		return fmt.Errorf("%w: collector service internal error: %s", domain.ErrInternal, st.Message())

	default:
		return fmt.Errorf("%w: unexpected error from collector: %s", domain.ErrInternal, st.Message())
	}
}
