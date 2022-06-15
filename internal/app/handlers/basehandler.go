package handlers

import (
	"net/http"

	"github.com/devkekops/wb-l0/internal/app/storage"
	"github.com/go-chi/chi/v5"
)

type BaseHandler struct {
	*chi.Mux
	orderRepo storage.OrderRepository
}

func NewBaseHandler(repo storage.OrderRepository) *BaseHandler {
	bh := &BaseHandler{
		Mux:       chi.NewMux(),
		orderRepo: repo,
	}

	bh.Get("/", bh.getOrders())

	return bh
}

func (bh *BaseHandler) getOrders() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

	}
}
