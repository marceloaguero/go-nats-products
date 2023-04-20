package delivery

import (
	"encoding/json"
	"log"

	"github.com/marceloaguero/go-nats-products/products/pkg/product"
	"github.com/nats-io/nats.go"
)

type delivery struct {
	usecase product.Usecase
	nc      *nats.Conn
}

func newDelivery(uc product.Usecase, nc *nats.Conn) *delivery {
	return &delivery{
		usecase: uc,
		nc:      nc,
	}
}

func NewDelivery(uc product.Usecase, natsURLs, subjPrefix, queue string) (*delivery, error) {
	nc, err := nats.Connect(natsURLs)
	if err != nil {
		return nil, err
	}

	delivery := newDelivery(uc, nc)

	err = Subscribe(delivery, nc, subjPrefix, queue)
	if err != nil {
		return nil, err
	}

	return delivery, nil
}

func (d *delivery) Drain() {
	d.nc.Drain()
}

func Subscribe(delivery *delivery, nc *nats.Conn, subjPrefix string, queue string) error {

	createSubj := subjPrefix + ".create"
	nc.QueueSubscribe(createSubj, queue, delivery.Create)

	return nil
}

// func (d *delivery) Create(subj, reply string, product *product.Product) {
func (d *delivery) Create(msg *nats.Msg) {
	product := &product.Product{}
	err := json.Unmarshal(msg.Data, &product)
	if err != nil {
		log.Println("Can't unmarshal product")
	}

	productCreated, err := d.usecase.Create(product)
	if err != nil {
		log.Printf("Create product failed. Err: %s", err.Error())
	}

	newProduct, err := json.Marshal(productCreated)
	if err != nil {
		log.Println("Can't marshal new product")
	}

	log.Printf("New product (unmarshaled): %v", productCreated)
	log.Printf("New product (marshaled): %v", newProduct)

	d.nc.Publish(msg.Reply, newProduct)
}
