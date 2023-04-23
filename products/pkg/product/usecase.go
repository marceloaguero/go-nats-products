package product

import (
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
)

var (
	validate *validator.Validate
)

// Usecase representa los casos de uso de productos
// Extiende la interface Product
type Usecase interface {
	Repository
	// UpdateStock permite modificar el stock de un producto
	// Por supuesto, se podría lograr el mismo resultado llamando al método Update de la interface Repository
	// pero lo agregamos sólo para extender dicha interface.
	UpdateStock(id uint, stock float64) (*Product, error)
}

type usecase struct {
	repository Repository
}

// NewUsecase creates a new usecase. Implements the Usecase interface
func NewUsecase(repo Repository) Usecase {
	return &usecase{
		repository: repo,
	}
}

// Create agrega un nuevo producto
func (u *usecase) Create(product *Product) (*Product, error) {
	// Verify name uniqueness
	product.Name = strings.TrimSpace(product.Name)
	_, err := u.GetByName(product.Name)
	if err == nil {
		return nil, errors.Errorf("UC - Create - Product with name %s already exists", product.Name)
	}

	validate := validator.New()
	err = validate.Struct(product)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return nil, errors.Wrap(validationErrors, "UC - Create - Error during product data validation")
	}

	product, err = u.repository.Create(product)
	if err != nil {
		return nil, errors.Wrap(err, "UC - Create - Error creating a new product")
	}

	return product, nil
}

// GetByID recupera un producto por ID
func (u *usecase) GetByID(id uint) (*Product, error) {
	product, err := u.repository.GetByID(id)
	if err != nil {
		return nil, errors.Wrap(err, "UC - GetByID - Error fetching a product")
	}

	return product, nil
}

// GetByName recupera un producto por su nombre
func (u *usecase) GetByName(name string) (*Product, error) {
	product, err := u.repository.GetByName(name)
	if err != nil {
		return nil, errors.Wrap(err, "UC - GetByName - Error fetching a product")
	}

	return product, nil
}

// GetAll recupera todos los productos
func (u *usecase) GetAll() ([]*Product, error) {
	products, err := u.repository.GetAll()
	if err != nil {
		return nil, errors.Wrap(err, "UC - GetAll - Error fetching all products")
	}

	return products, nil
}

// Update modifica un producto existente
func (u *usecase) Update(product *Product) (*Product, error) {
	// Trim spaces
	product.Name = strings.TrimSpace(product.Name)

	formerProduct := &Product{}

	// Verificar la unicidad del nombre
	formerProduct, err := u.GetByName(product.Name)
	if (err == nil) && (formerProduct.ID != product.ID) {
		return nil, errors.Errorf("UC - Update - Product with name %s already exists", product.Name)
	}

	validate = validator.New()
	if err := validate.Struct(product); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return nil, errors.Wrap(validationErrors, "UC - Update - Error during product data validation")
	}

	formerProduct, err = u.GetByID(product.ID)
	if err != nil {
		return nil, errors.Wrapf(err, "UC - Update - Product with id %d does not exist", product.ID)
	}

	product, err = u.repository.Update(product)
	if err != nil {
		return nil, errors.Wrapf(err, "UC - Update - Error updating product with id %d", product.ID)
	}

	return product, nil
}

// Delete elimina un producto
func (u *usecase) Delete(product *Product) error {
	_, err := u.GetByID(product.ID)
	if err != nil {
		return errors.Wrapf(err, "UC - Delete - Product with id %d does not exist", product.ID)
	}

	if err := u.repository.Delete(product); err != nil {
		return errors.Wrapf(err, "UC - Delete - Error deleting product with id %d", product.ID)
	}

	return nil
}

// UpdateStock permite modificar el stock de un producto
func (u *usecase) UpdateStock(id uint, stock float64) (*Product, error) {
	product, err := u.GetByID(id)
	if err != nil {
		return nil, errors.Wrapf(err, "UC - UpdateStock - Product with id %d does not exist", id)
	}

	if stock < 0 {
		return nil, errors.New("UC - UpdateStock - Stock can't be negative")
	}

	product.Stock = stock
	_, err = u.Update(product)

	return product, err
}
