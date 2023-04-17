package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/marceloaguero/go-nats-products/gateway/pkg/delivery/products"
	"github.com/marceloaguero/go-nats-products/gateway/pkg/delivery/router"
	"github.com/nats-io/nats.go"
)

func main() {
	pathPrefix := os.Getenv("PATH_PREFIX")
	natsURLs := os.Getenv("NATS_URLS")
	productsSubjPrefix := os.Getenv("PRODUCTS_SUBJ_PREFIX")
	productsQueue := os.Getenv("PRODUCTS_QUEUE")

	// Connect to NATS server
	nc, err := nats.Connect(natsURLs)
	if err != nil {
		log.Panic(err)
	}
	defer nc.Close()

	ec, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		log.Panic(err)
	}
	defer ec.Close()

	productsDelivery := products.NewDelivery(ec, productsSubjPrefix, productsQueue)

	_, err = router.NewRouter(productsDelivery, pathPrefix)
	if err != nil {
		log.Panic(err)
	}

	// Setup an interrupt handler to drain nats
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	log.Println()
	log.Printf("Draining...")
	nc.Drain()
	log.Fatal("Exiting...")
}
