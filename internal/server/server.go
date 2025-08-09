package server

import (
	"net/http"

	"github.com/ghostchain1/core-service/internal/config"
	"github.com/ghostchain1/core-service/pkg/version"
	"github.com/go-chi/chi/v5"
)

func New(cfg *config.Config) *http.Server {
	r := chi.NewRouter()

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	r.Get("/readyz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ready"))
	})

	r.Get("/version", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(version.Version))
	})

	return &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}
}
