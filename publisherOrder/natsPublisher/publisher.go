package natsPublisher

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/nats-io/stan.go"

	"awesomeProject5/publisherOrder/modelsNats"
)

var (
	clusterID = "test-cluster"
	clientID  = "test-client2"
	channel   = "test-channel"
)

type Publisher struct {
	s *chi.Mux
	c stan.Conn
}

func NewPublisher() *Publisher {
	conn, err := stan.Connect(clusterID, clientID)
	if err != nil {
		log.Fatal(err)
	}
	return &Publisher{
		s: chi.NewRouter(),
		c: conn,
	}
}

func (p *Publisher) Run() error {

	if err := p.configRouting(); err != nil {
		log.Fatal(err)

	}

	http.ListenAndServe("localhost:8081", p.s)

	return nil

}

func (p *Publisher) configRouting() error {
	r := p.s

	r.Post("/postorder", p.SaveProducts())

	return nil
}

func (p *Publisher) SaveProducts() http.HandlerFunc {

	return func(writer http.ResponseWriter, request *http.Request) {
		orderReq := &modelsNats.Order{}
		if err := json.NewDecoder(request.Body).Decode(orderReq); err != nil {
			http.Error(writer, "invalid json", http.StatusBadRequest)
			return
		}

		order, err := json.Marshal(orderReq)
		fmt.Println(order)
		if err != nil {
			log.Fatal(err)
		}
		err = p.c.Publish(channel, order)
	}
}
