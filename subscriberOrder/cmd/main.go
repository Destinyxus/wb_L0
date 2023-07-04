package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"awesomeProject5/subscriberOrder/internal/apiserver"
	"awesomeProject5/subscriberOrder/internal/nats"
	store2 "awesomeProject5/subscriberOrder/internal/store"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	store := store2.NewStore()
	cache := store2.NewCache()
	s := apiserver.NewAPIServer(store, cache)
	sub := nats.NewSubscriber(store, cache)

	server := &http.Server{
		Addr:    ":8080",
		Handler: s.Router,
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		if err := s.Init(); err != nil {
			panic(err)
		}
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal(err)
			return
		}
	}()

	ctxWT, cancel2 := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel2()

	go func() {
		select {
		case <-sigChan:
			fmt.Println("Graceful shutdown started...")
			store.CloseConnectionDB()
			fmt.Println("DB connection closed")
			fmt.Println("Subscription closed")
			cancel()
			server.Shutdown(ctxWT)
			fmt.Println("API server closed")
			return
		}

	}()

	time.Sleep(time.Second * 2)

	sub.InitSubscriberConn(ctx)
	wg.Wait()
}
