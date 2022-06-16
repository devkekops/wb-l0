package main

import (
	"flag"
	"log"

	"github.com/caarlos0/env"
	"github.com/devkekops/wb-l0/internal/app/config"
	"github.com/devkekops/wb-l0/internal/app/server"
)

func main() {
	cfg := config.Config{
		RunAddress:  "127.0.0.1:8080",
		DatabaseURI: "postgres://localhost:5432/wb_l0_orders",
		NatsURI:     "nats://127.0.0.1:4222",
		NatsSubject: "hello",
	}

	if err := env.Parse(&cfg); err != nil {
		log.Fatal(err)
		return
	}

	flag.StringVar(&cfg.RunAddress, "a", cfg.RunAddress, "run address")
	flag.StringVar(&cfg.DatabaseURI, "d", cfg.DatabaseURI, "database URI")
	flag.StringVar(&cfg.NatsURI, "n", cfg.NatsURI, "nats URI")
	flag.StringVar(&cfg.NatsURI, "s", cfg.NatsURI, "nats subject")
	flag.Parse()

	log.Fatal(server.Serve(&cfg))
}
