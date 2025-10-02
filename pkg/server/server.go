package server

import (
	"go-pismo-challenge/pkg/handler"
	"net/http"
)

func NewServer(a *handler.Account, t *handler.Transaction) *http.Server {
	r := NewRouter(a, t)

	return &http.Server{
		Addr:    ":3000",
		Handler: r,
	}
}
