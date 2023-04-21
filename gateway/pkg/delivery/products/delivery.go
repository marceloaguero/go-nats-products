package products

import (
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nats-io/nats.go"
)

const (
	timeout = time.Millisecond * 500
)

type Delivery interface {
	Create(c *gin.Context)
}

type delivery struct {
	nc         *nats.Conn
	subjPrefix string
	queue      string
}

func NewDelivery(nc *nats.Conn, subjPrefix, queue string) Delivery {
	return &delivery{
		nc:         nc,
		subjPrefix: subjPrefix,
		queue:      queue,
	}
}

func (d *delivery) Create(c *gin.Context) {
	createSubj := d.subjPrefix + ".create"
	data, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
	}

	msg, err := d.nc.Request(createSubj, data, timeout)
	if err != nil {
		log.Printf("Request err: %v", err)
	}

	c.Data(http.StatusCreated, "application/json", msg.Data)
}
