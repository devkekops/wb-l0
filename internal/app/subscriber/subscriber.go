package subscriber

import (
	"fmt"
	"runtime"

	"github.com/devkekops/wb-l0/internal/app/storage"
	"github.com/nats-io/nats.go"
)

type Subscriber struct {
	c *nats.EncodedConn
}

func NewSubscriber() *Subscriber {
	nc, _ := nats.Connect(nats.DefaultURL)
	c, _ := nats.NewEncodedConn(nc, nats.JSON_ENCODER)

	return &Subscriber{
		c: c,
	}
}

func (s *Subscriber) Check() {
	s.c.Subscribe("hello", func(order *storage.Order) {
		fmt.Printf("Received a order: %+v\n", order)
	})
	s.c.Flush()

	runtime.Goexit()
}
