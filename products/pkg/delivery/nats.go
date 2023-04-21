package delivery

import (
	"encoding/json"
	"log"

	"clevergo.tech/jsend"
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

func (d *delivery) Create(msg *nats.Msg) {
	product := &product.Product{}
	err := json.Unmarshal(msg.Data, &product)
	if err != nil {
		log.Println("Can't unmarshal product")
		jsendReply := jsend.NewFail(err.Error())
		reply, _ := json.Marshal(&jsendReply)

		d.nc.Publish(msg.Reply, reply)
		return
	}

	productCreated, err := d.usecase.Create(product)
	if err != nil {
		log.Printf("Create product failed. Err: %s", err.Error())
		jsendReply := jsend.NewFail(err.Error())
		reply, err := json.Marshal(&jsendReply)
		if err != nil {
			log.Println("Can't marshal jsend reply")
		}

		d.nc.Publish(msg.Reply, reply)
		return
	}

	jsendReply := jsend.New(productCreated)
	reply, err := json.Marshal(&jsendReply)
	if err != nil {
		log.Println("Can't marshal jsend reply")
	}

	d.nc.Publish(msg.Reply, reply)
}
