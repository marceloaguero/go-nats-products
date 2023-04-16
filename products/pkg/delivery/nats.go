package delivery

import (
	"github.com/marceloaguero/go-nats-products/products/pkg/product"
	"github.com/nats-io/nats.go"
)

type delivery struct {
	usecase product.Usecase
	ec      *nats.EncodedConn
}

func newDelivery(uc product.Usecase, ec *nats.EncodedConn) *delivery {
	return &delivery{
		usecase: uc,
		ec:      ec,
	}
}

func NewDelivery(uc product.Usecase, natsURLs, subjPrefix, queue string) (*delivery, error) {
	nc, err := nats.Connect(natsURLs)
	if err != nil {
		return nil, err
	}

	ec, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		return nil, err
	}

	delivery := newDelivery(uc, ec)

	err = Subscribe(delivery, ec, subjPrefix, queue)
	if err != nil {
		return nil, err
	}

	return delivery, nil
}

func (d *delivery) Drain() {
	d.ec.Drain()
}

func Subscribe(delivery *delivery, ec *nats.EncodedConn, subjPrefix string, queue string) error {

	createSubj := subjPrefix + ".create"
	ec.QueueSubscribe(createSubj, queue, delivery.Create)

	return nil
}

func (d *delivery) Create(subj, reply string, product *product.Product) {
	productCreated, _ := d.usecase.Create(product)
	//newProduct, _ := json.Marshal(productCreated)

	//log.Printf("Product: %v", newProduct)
	//d.ec.Publish(reply, newProduct)
	d.ec.Publish(reply, productCreated)
}
