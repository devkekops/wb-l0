package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/nats-io/nats.go"
)

func main() {
	nc, _ := nats.Connect(nats.DefaultURL)

	file, err := os.OpenFile("cmd/publisher/model.json", os.O_RDONLY, 0777)
	defer func() {
		if err = file.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	nc.Publish("hello", b)
	nc.Flush()
}
