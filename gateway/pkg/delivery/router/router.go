package router

import (
	"github.com/gin-gonic/gin"

	"github.com/marceloaguero/go-nats-products/gateway/pkg/delivery/products"
)

type router struct {
	productsDelivery products.Delivery
}

func NewRouter(productsDelivery products.Delivery, pathPrefix string) (*router, error) {
	router := &router{
		productsDelivery: productsDelivery,
	}

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	products := r.Group("/products")
	{
		products.GET("/", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "productos",
			})
		})
		// Crear un nuevo producto
		products.POST("/", router.productsDelivery.Create)
		products.GET("/:id", router.productsDelivery.GetByID)
	}

	err := r.Run()
	if err != nil {
		return nil, err
	}

	return router, nil
}
