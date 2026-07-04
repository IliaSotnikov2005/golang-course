package v1

import (
	"context"
	"log/slog"
	"net"
	"net/http"

	"github.com/IliaSotnikov2005/golang-course/task6/repo-stat/api/internal/controller/http/respond"
	"github.com/IliaSotnikov2005/golang-course/task6/repo-stat/api/internal/domain"
)

type Limiter interface {
	Allow(ctx context.Context, ip string, rps float64, burst int) (bool, error)
}

func RateLimit(log *slog.Logger, limiter Limiter, rps float64, burst int) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				ip = r.RemoteAddr
			}

			allowed, err := limiter.Allow(r.Context(), ip, rps, burst)
			if err != nil {
				log.Error("rate limit error", slog.Any("error", err))
				next.ServeHTTP(w, r)
				return
			}

			if !allowed {
				log.Warn("rate limit exceeded", slog.String("ip", ip))
				respond.Error(w, domain.ErrRateLimit)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
