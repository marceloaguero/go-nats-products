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
	var s string

	s = subjPrefix + ".create"
	_, err := nc.QueueSubscribe(s, queue, delivery.Create)

	s = subjPrefix + ".getbyid"
	_, err = nc.QueueSubscribe(s, queue, delivery.GetByID)
	/*
		s = subjPrefix + ".getbyname"
		_, err = nc.QueueSubscribe(s, queue, delivery.GetByName)

		s = subjPrefix + ".getall"
		_, err = nc.QueueSubscribe(s, queue, delivery.GetAll)

		s = subjPrefix + ".update"
		_, err = nc.QueueSubscribe(s, queue, delivery.Update)

		s = subjPrefix + ".delete"
		_, err = nc.QueueSubscribe(s, queue, delivery.Delete)

		s = subjPrefix + ".updatestock"
		_, err = nc.QueueSubscribe(s, queue, delivery.UpdateStock)
	*/
	return err
}

func (d *delivery) Create(msg *nats.Msg) {
	product := &product.Product{}
	err := json.Unmarshal(msg.Data, &product)
	if err != nil {
		jsendFailReply := jsend.NewFail(err.Error())
		failReply, _ := json.Marshal(&jsendFailReply)

		d.nc.Publish(msg.Reply, failReply)
		return
	}

	productCreated, err := d.usecase.Create(product)
	if err != nil {
		jsendFailReply := jsend.NewFail(err.Error())
		failReply, _ := json.Marshal(&jsendFailReply)

		d.nc.Publish(msg.Reply, failReply)
		return
	}

	jsendReply := jsend.New(productCreated)
	reply, err := json.Marshal(&jsendReply)
	if err != nil {
		log.Println("DLV - Create - Can't marshal jsend reply")
	}

	d.nc.Publish(msg.Reply, reply)
}

func (d *delivery) GetByID(msg *nats.Msg) {
	product := &product.Product{}
	err := json.Unmarshal(msg.Data, &product)
	if err != nil {
		jsendFailReply := jsend.NewFail(err.Error())
		failReply, _ := json.Marshal(&jsendFailReply)

		d.nc.Publish(msg.Reply, failReply)
		return
	}

	productRetrieved, err := d.usecase.GetByID(product.ID)
	if err != nil {
		jsendFailReply := jsend.NewFail(err.Error())
		failReply, _ := json.Marshal(&jsendFailReply)

		d.nc.Publish(msg.Reply, failReply)
		return
	}

	jsendReply := jsend.New(productRetrieved)
	reply, err := json.Marshal(&jsendReply)
	if err != nil {
		log.Println("DLV - GetByID - Can't marshal jsend reply")
	}

	d.nc.Publish(msg.Reply, reply)
}
