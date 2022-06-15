package main

import (
	"log"

	"github.com/devkekops/wb-l0/internal/app/server"
)

func main() {
	log.Fatal(server.Serve())
}
