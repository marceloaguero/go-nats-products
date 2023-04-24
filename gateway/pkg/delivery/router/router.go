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
		// Crear un nuevo producto
		products.POST("/", router.productsDelivery.Create)
		// Recuperar todos los productos
		products.GET("/", router.productsDelivery.GetAll)
		// Recuperar un producto por su ID
		products.GET("/:id", router.productsDelivery.GetByID)
		// Recuperar producto por nombre
		products.GET("/names/:name", router.productsDelivery.GetByName)
		// Modificar un producto
		products.PUT("/:id", router.productsDelivery.Update)
		// Eliminar un producto
		products.DELETE("/:id", router.productsDelivery.Delete)
	}

	err := r.Run()
	if err != nil {
		return nil, err
	}

	return router, nil
}
