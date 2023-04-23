package products

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/marceloaguero/go-nats-products/gateway/pkg/delivery/jsenderrors"
	"github.com/nats-io/nats.go"
)

const (
	timeout = time.Millisecond * 500
)

type NatsMsgData struct {
	Status string `json:"status"`
}

type Delivery interface {
	Create(c *gin.Context)
	GetByID(c *gin.Context)
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

func sendRequest(c *gin.Context, d *delivery, method string, subj string, request []byte, defaultStatus int) {
	msg, err := d.nc.Request(subj, request, timeout)
	if err != nil {
		log.Printf("%s - Request error: %s", method, err.Error())
		jsenderrors.ReturnError(c, err.Error())
		return
	}

	msgData := &NatsMsgData{}
	err = json.Unmarshal(msg.Data, &msgData)
	if err != nil {
		log.Printf("%s - Unmarshal reply error: %s", method, err.Error())
		jsenderrors.ReturnError(c, err.Error())
		return
	}

	stat := msgData.Status
	var httpStatus int
	switch {
	case stat == "fail":
		httpStatus = http.StatusBadRequest
		break
	case stat == "error":
		httpStatus = http.StatusInternalServerError
		break
	default:
		httpStatus = defaultStatus
	}
	c.Data(httpStatus, "application/json", msg.Data)
}

func (d *delivery) Create(c *gin.Context) {
	subj := d.subjPrefix + ".create"
	request, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		jsenderrors.ReturnError(c, err.Error())
		return
	}

	method := "DLV - Products - Create"
	sendRequest(c, d, method, subj, request, http.StatusCreated)
}

func (d *delivery) GetByID(c *gin.Context) {
	id := c.Param("id")
	request := []byte(fmt.Sprintf("{ \"id\": %s }", id))
	subj := d.subjPrefix + ".getbyid"

	method := "DLV - Products - GetByID"
	sendRequest(c, d, method, subj, request, http.StatusOK)
}
