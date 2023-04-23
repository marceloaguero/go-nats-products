package mysql_orm

import (
	"fmt"

	"github.com/marceloaguero/go-nats-products/products/pkg/product"
	"github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type ormRepo struct {
	db *gorm.DB
}

// NewRepo crea un repositorio implementado en ORM (MySQL)
func NewRepo(dsName, dbName string) (product.Repository, error) {
	db, err := dbConnect(dsName, dbName)
	if err != nil {
		return nil, errors.Wrap(err, "MySQL ORM - Can't connect to DB")
	}

	db.AutoMigrate(&product.Product{})

	return &ormRepo{
		db: db,
	}, nil
}

func dbConnect(dsName, dbName string) (*gorm.DB, error) {
	conn := fmt.Sprintf("%s/%s?charset=utf8&parseTime=True&loc=Local", dsName, dbName)

	db, err := gorm.Open(mysql.Open(conn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (r *ormRepo) Create(product *product.Product) (*product.Product, error) {
	result := r.db.Create(&product)
	return product, result.Error
}

func (r *ormRepo) GetByID(id uint) (*product.Product, error) {
	var product product.Product
	result := r.db.Take(&product, id)
	return &product, result.Error
}

func (r *ormRepo) GetByName(name string) (*product.Product, error) {
	var product product.Product
	result := r.db.Take(&product, "name = ?", name)
	return &product, result.Error
}

func (r *ormRepo) GetAll() ([]*product.Product, error) {
	products := []*product.Product{}
	result := r.db.Find(&products)
	return products, result.Error
}

func (r *ormRepo) Update(product *product.Product) (*product.Product, error) {
	result := r.db.Model(&product).Updates(product)
	if result.Error == nil {
		r.db.Model(&product).Updates(map[string]interface{}{"is_active": product.IsActive})
	}
	return product, result.Error
}

func (r *ormRepo) Delete(product *product.Product) error {
	result := r.db.Delete(&product)
	return result.Error
}
