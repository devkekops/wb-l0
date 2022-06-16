package subscriber

import (
	"log"

	"github.com/devkekops/wb-l0/internal/app/storage"
	"github.com/nats-io/nats.go"
)

type Subscriber struct {
	c      *nats.EncodedConn
	r      storage.OrderRepository
	recvCh chan *storage.Order
}

func NewSubscriber(natsURI string, natsSubject string, repo storage.OrderRepository) (*Subscriber, error) {
	nc, err := nats.Connect(natsURI)
	if err != nil {
		return nil, err
	}
	c, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		return nil, err
	}

	recvCh := make(chan *storage.Order)
	c.BindRecvChan(natsSubject, recvCh)

	return &Subscriber{
		c:      c,
		r:      repo,
		recvCh: recvCh,
	}, nil
}

func (s *Subscriber) Check() {
	for {
		order := <-s.recvCh
		//fmt.Printf("Received a order: %+v\n", order)
		err := s.r.SaveOrder(*order)
		if err != nil {
			log.Println(err)
		}
	}
}
