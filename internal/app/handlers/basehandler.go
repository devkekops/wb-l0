package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/devkekops/wb-l0/internal/app/storage"
	"github.com/go-chi/chi/middleware"
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

	bh.Use(middleware.Logger)
	bh.Get("/{id}", bh.getOrders())

	return bh
}

func (bh *BaseHandler) getOrders() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		orderID := chi.URLParam(req, "id")
		order, err := bh.orderRepo.GetOrderByID(orderID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}

		buf, err := json.Marshal(order)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(buf)
		if err != nil {
			log.Println(err)
		}
	}
}
