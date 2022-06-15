package server

import (
	"log"
	"net/http"

	"github.com/devkekops/wb-l0/internal/app/config"
	"github.com/devkekops/wb-l0/internal/app/handlers"
	"github.com/devkekops/wb-l0/internal/app/storage"
	"github.com/devkekops/wb-l0/internal/app/subscriber"
)

func Serve(cfg *config.Config) error {
	orderRepo, err := storage.NewOrderRepoDB(cfg.DatabaseURI)
	if err != nil {
		log.Fatal(err)
	}

	subscriber, err := subscriber.NewSubscriber(cfg.NatsURI, orderRepo)
	if err != nil {
		log.Fatal(err)
	}
	go subscriber.Check()

	baseHandler := handlers.NewBaseHandler(orderRepo)

	server := &http.Server{
		Addr:    cfg.RunAddress,
		Handler: baseHandler,
	}

	return server.ListenAndServe()
}
