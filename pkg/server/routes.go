package server

import (
	_ "go-pismo-challenge/docs"
	"go-pismo-challenge/pkg/handler"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

func NewRouter(a *handler.Account, t *handler.Transaction) *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	// accounts
	router.Route("/v1/accounts", func(r chi.Router) {
		r.Post("/", a.Create)
		r.Get("/{uuid}", a.Get)
	})

	// transactions
	router.Route("/v1/transactions", func(r chi.Router) {
		r.Post("/", t.Create)
	})

	// serve swagger UI
	router.Get("/swagger/*", httpSwagger.WrapHandler)

	return router
}
