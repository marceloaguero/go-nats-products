package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/marceloaguero/go-nats-products/products/pkg/delivery"
	"github.com/marceloaguero/go-nats-products/products/pkg/product"
	repo "github.com/marceloaguero/go-nats-products/products/pkg/repository"
)

func main() {
	dbDsn := os.Getenv("DB_DSN")
	dbName := os.Getenv("DB_NAME")
	natsURLs := os.Getenv("NATS_URLS")
	subjPrefix := os.Getenv("SUBJ_PREFIX")
	queue := os.Getenv("QUEUE")

	repository, err := repo.NewRepo(dbDsn, dbName)
	if err != nil {
		log.Panic(err)
	}

	usecase := product.NewUsecase(repository)

	delivery, err := delivery.NewDelivery(usecase, natsURLs, subjPrefix, queue)
	if err != nil {
		log.Panic(err)
	}

	// Setup the interrupt handler to drain so we don't miss
	// requests when scaling down.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	log.Printf("Draining...")
	delivery.Drain()
	log.Fatalf("Exiting")
}
