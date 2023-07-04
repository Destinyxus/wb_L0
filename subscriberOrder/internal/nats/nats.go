package nats

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/nats-io/stan.go"

	"awesomeProject5/subscriberOrder/internal/models"
	"awesomeProject5/subscriberOrder/internal/store"
	"awesomeProject5/subscriberOrder/utils/validation"
)

type Subscriber struct {
	store *store.Store
	cache *store.Cache
}

func NewSubscriber(store *store.Store, cache *store.Cache) *Subscriber {
	return &Subscriber{
		store: store,
		cache: cache,
	}
}

var (
	clusterID = "test-cluster"
	clientID  = "test-client"
	channel   = "test-channel"
)

func (s *Subscriber) InitSubscriberConn(ctx context.Context) {

	fmt.Println("here")
	sc, err := stan.Connect(clusterID, clientID)
	if err != nil {
		log.Fatal(err)
	}
	defer func(sc stan.Conn) {
		err := sc.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(sc)
	durableName := "test-Name"
	sub, err := sc.Subscribe(channel, func(msg *stan.Msg) {

		orderReq, err := validation.Validate(msg.Data)
		if err != nil {
			fmt.Println(err)

		} else {
			err = s.store.AddOrder(orderReq)
			err = s.store.AddDelivery(orderReq.Delivery)
			err = s.store.AddPayments(orderReq.Payment)
			err = s.store.AddItems(orderReq.Items)
			err = s.cache.SaveToCache(msg.Data, strconv.Itoa(models.OrderID))
		}

	}, stan.DurableName(durableName))
	select {
	case <-ctx.Done():
		err := sub.Close()
		if err != nil {
			return
		}
		err = sc.Close()
		if err != nil {
			return
		}

	}
}
