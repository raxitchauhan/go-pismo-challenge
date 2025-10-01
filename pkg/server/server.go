package server

import (
	"go-pismo-challenge/pkg/handler"
	"net/http"
)

func Start(a *handler.Account, t *handler.Transaction) error {
	r := NewRouter(a, t)

	return http.ListenAndServe(":3000", r)
}
