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

	return err
}

func JsendFailReply(d *delivery, msg *nats.Msg, errMsg string) {
	jsendFailReply := jsend.NewFail(errMsg)
	failReply, _ := json.Marshal(&jsendFailReply)

	d.nc.Publish(msg.Reply, failReply)
}

func (d *delivery) Create(msg *nats.Msg) {
	product := &product.Product{}
	err := json.Unmarshal(msg.Data, &product)
	if err != nil {
		JsendFailReply(d, msg, err.Error())
		return
	}

	productCreated, err := d.usecase.Create(product)
	if err != nil {
		JsendFailReply(d, msg, err.Error())
		return
	}

	jsendReply := jsend.New(productCreated)
	reply, err := json.Marshal(&jsendReply)
	if err != nil {
		log.Println("DLV - Create - Can't marshal jsend reply")
		JsendFailReply(d, msg, err.Error())
		return
	}

	d.nc.Publish(msg.Reply, reply)
}

func (d *delivery) GetByID(msg *nats.Msg) {
	product := &product.Product{}
	err := json.Unmarshal(msg.Data, &product)
	if err != nil {
		JsendFailReply(d, msg, err.Error())
		return
	}

	productRetrieved, err := d.usecase.GetByID(product.ID)
	if err != nil {
		JsendFailReply(d, msg, err.Error())
		return
	}

	jsendReply := jsend.New(productRetrieved)
	reply, err := json.Marshal(&jsendReply)
	if err != nil {
		log.Println("DLV - GetByID - Can't marshal jsend reply")
		JsendFailReply(d, msg, err.Error())
		return
	}

	d.nc.Publish(msg.Reply, reply)
}

func (d *delivery) GetByName(msg *nats.Msg) {
	product := &product.Product{}
	err := json.Unmarshal(msg.Data, &product)
	if err != nil {
		JsendFailReply(d, msg, err.Error())
		return
	}

	productRetrieved, err := d.usecase.GetByName(product.Name)
	if err != nil {
		JsendFailReply(d, msg, err.Error())
		return
	}

	jsendReply := jsend.New(productRetrieved)
	reply, err := json.Marshal(&jsendReply)
	if err != nil {
		log.Println("DLV - GetByName - Can't marshal jsend reply")
		JsendFailReply(d, msg, err.Error())
		return
	}

	d.nc.Publish(msg.Reply, reply)
}

func (d *delivery) GetAll(msg *nats.Msg) {
	products, err := d.usecase.GetAll()
	if err != nil {
		JsendFailReply(d, msg, err.Error())
		return
	}

	jsendReply := jsend.New(products)
	reply, err := json.Marshal(&jsendReply)
	if err != nil {
		log.Println("DLV - GetAll - Can't marshal jsend reply")
		JsendFailReply(d, msg, err.Error())
		return
	}

	d.nc.Publish(msg.Reply, reply)
}

func (d *delivery) Update(msg *nats.Msg) {
	product := &product.Product{}
	err := json.Unmarshal(msg.Data, &product)
	if err != nil {
		JsendFailReply(d, msg, err.Error())
		return
	}

	productUpdated, err := d.usecase.Update(product)
	if err != nil {
		JsendFailReply(d, msg, err.Error())
		return
	}

	jsendReply := jsend.New(productUpdated)
	reply, err := json.Marshal(&jsendReply)
	if err != nil {
		log.Println("DLV - Update - Can't marshal jsend reply")
		JsendFailReply(d, msg, err.Error())
		return
	}

	d.nc.Publish(msg.Reply, reply)
}

func (d *delivery) Delete(msg *nats.Msg) {
	product := &product.Product{}
	err := json.Unmarshal(msg.Data, &product)
	if err != nil {
		JsendFailReply(d, msg, err.Error())
		return
	}

	err = d.usecase.Delete(product)
	if err != nil {
		JsendFailReply(d, msg, err.Error())
		return
	}

	jsendReply := jsend.New(nil)
	reply, err := json.Marshal(&jsendReply)
	if err != nil {
		log.Println("DLV - Delete - Can't marshal jsend reply")
		JsendFailReply(d, msg, err.Error())
		return
	}

	d.nc.Publish(msg.Reply, reply)
}

func (d *delivery) UpdateStock(msg *nats.Msg) {
	product := &product.Product{}
	err := json.Unmarshal(msg.Data, &product)
	if err != nil {
		JsendFailReply(d, msg, err.Error())
		return
	}

	productUpdated, err := d.usecase.UpdateStock(product.ID, product.Stock)
	if err != nil {
		JsendFailReply(d, msg, err.Error())
		return
	}

	jsendReply := jsend.New(productUpdated)
	reply, err := json.Marshal(&jsendReply)
	if err != nil {
		log.Println("DLV - UpdateStock - Can't marshal jsend reply")
		JsendFailReply(d, msg, err.Error())
		return
	}

	d.nc.Publish(msg.Reply, reply)
}
