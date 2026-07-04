package v1

import (
	"log/slog"
	"net/http"
	"time"

	redismiddleware "github.com/IliaSotnikov2005/golang-course/task6/repo-stat/api/internal/controller/http/v1/middleware"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger"
)

func slogRequestLogger(log *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			next.ServeHTTP(ww, r)
			log.Info("request",
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.Int("status", ww.Status()),
				slog.Duration("duration", time.Since(start)),
			)
		})
	}
}

func (h *Handler) Router(limiter redismiddleware.Limiter, rps float64, burst int) chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(slogRequestLogger(h.log))
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodOptions},
		AllowedHeaders: []string{"Accept"},
		MaxAge:         300,
	}))

	r.Get("/swagger/*", h.serveSwagger)

	r.Route("/api/v1", func(r chi.Router) {
		r.Use(redismiddleware.RateLimit(h.log, limiter, rps, burst))

		r.Get("/repositories/info", h.getRepository)

		r.Route("/subscriptions", func(r chi.Router) {
			r.Get("/", h.listSubscriptions)
			r.Get("/info", h.getSubscriptionsInfo)

			r.Route("/{owner}/{repo}", func(r chi.Router) {
				r.Post("/", h.subscribe)
				r.Delete("/", h.unsubscribe)
			})
		})
	})

	r.Route("/api", func(r chi.Router) {
		r.Get("/ping", h.healthCheck)
	})

	return r
}

func (h *Handler) serveSwagger(w http.ResponseWriter, r *http.Request) {
	httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	).ServeHTTP(w, r)
}
