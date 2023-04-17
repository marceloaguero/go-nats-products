package products

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/marceloaguero/go-nats-products/products/pkg/product"
	"github.com/nats-io/nats.go"
)

const (
	timeout = time.Millisecond * 500
)

type Delivery interface {
	Create(c *gin.Context)
}

type delivery struct {
	ec         *nats.EncodedConn
	subjPrefix string
	queue      string
}

func NewDelivery(ec *nats.EncodedConn, subjPrefix, queue string) Delivery {
	return &delivery{
		ec:         ec,
		subjPrefix: subjPrefix,
		queue:      queue,
	}
}

func (d *delivery) Create(c *gin.Context) {
	var newProduct *product.Product
	productCreated := &product.Product{}
	createSubj := d.subjPrefix + ".create"
	//data, err := ioutil.ReadAll(c.Request.Body)
	err := c.BindJSON(&newProduct)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
	}
	err = d.ec.Request(createSubj, newProduct, productCreated, timeout)
	if err != nil {
		log.Printf("err: %v", err)
	}
	c.IndentedJSON(http.StatusOK, productCreated)
}
