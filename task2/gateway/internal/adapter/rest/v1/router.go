package v1

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/IliaSotnikov2005/golang-course/task2/gateway/docs"
)

func (h *Handler) Router() chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{http.MethodGet},
		AllowedHeaders: []string{"Accept"},
		MaxAge:         300,
	}))

	r.Get("/swagger/*", h.serveSwagger)

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/repos/{owner}/{repo}", h.getRepository)
	})

	r.Route("/api", func(r chi.Router) {
		r.Get("/health", h.healthCheck)
	})

	return r
}

func (h *Handler) serveSwagger(w http.ResponseWriter, r *http.Request) {
	httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	).ServeHTTP(w, r)
}
