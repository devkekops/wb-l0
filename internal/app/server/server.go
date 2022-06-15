package server

import (
	"net/http"

	"github.com/devkekops/wb-l0/internal/app/handlers"
	"github.com/devkekops/wb-l0/internal/app/storage"
	"github.com/devkekops/wb-l0/internal/app/subscriber"
)

func Serve() error {
	orderRepo := storage.NewOderRepo()
	baseHandler := handlers.NewBaseHandler(orderRepo)
	subscriber := subscriber.NewSubscriber()
	subscriber.Check()

	server := &http.Server{
		Addr:    "127.0.0.1:8080",
		Handler: baseHandler,
	}

	return server.ListenAndServe()
}
