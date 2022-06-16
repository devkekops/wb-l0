package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/nats-io/nats.go"
)

func main() {
	nc, _ := nats.Connect(nats.DefaultURL)

	dir := "cmd/publisher/models/"

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		path := dir + f.Name()
		if filepath.Ext(path) == ".json" {
			fmt.Println(path)
			file, err := os.OpenFile(path, os.O_RDONLY, 0777)
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
	}

}
