package subscriber

import (
	"fmt"
	"log"
	"runtime"

	"github.com/devkekops/wb-l0/internal/app/storage"
	"github.com/nats-io/nats.go"
)

type Subscriber struct {
	c *nats.EncodedConn
	r storage.OrderRepository
}

func NewSubscriber(natsURI string, repo storage.OrderRepository) (*Subscriber, error) {
	nc, err := nats.Connect(natsURI)
	if err != nil {
		return nil, err
	}
	c, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		return nil, err
	}

	return &Subscriber{
		c: c,
		r: repo,
	}, nil
}

func (s *Subscriber) Check() {
	s.c.Subscribe("hello", func(order *storage.Order) {
		fmt.Printf("Received a order: %+v\n", order)
		err := s.r.SaveOrder(*order)
		if err != nil {
			log.Println(err)
		}
	})
	s.c.Flush()

	runtime.Goexit()
}
