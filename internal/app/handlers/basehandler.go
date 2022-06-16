package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/devkekops/wb-l0/internal/app/storage"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

type BaseHandler struct {
	*chi.Mux
	orderRepo storage.OrderRepository
	fs        http.Handler
}

func NewBaseHandler(repo storage.OrderRepository) *BaseHandler {
	root := "./internal/app/static"
	fs := http.FileServer(http.Dir(root))

	bh := &BaseHandler{
		Mux:       chi.NewMux(),
		orderRepo: repo,
		fs:        fs,
	}
	bh.Use(middleware.Logger)

	bh.Get("/", bh.getIndex())
	bh.Get("/{id}", bh.getOrder())

	return bh
}

func (bh *BaseHandler) getIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bh.fs.ServeHTTP(w, r)
	}
}

func (bh *BaseHandler) getOrder() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		orderID := chi.URLParam(req, "id")
		order, err := bh.orderRepo.GetOrderByID(orderID)
		if err != nil {
			if errors.Is(err, storage.ErrOrderNotExists) {
				log.Println(err)
				w.WriteHeader(http.StatusNoContent)
				return
			}
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
