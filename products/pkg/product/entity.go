package product

// Product describe un producto en el sistema
// Se utiliza gorm (https://gorm.io/) para modelar la entidad Product en la base de datos.
// Es por ello que en la declaración de los atributos, además del nombre que recibe el atributo en json,
// también se declaran las características necesarias de gorm.
type Product struct {
	ID          uint    `json:"id" gorm:"primaryKey"`                                     // Identificador del producto. Es clave primaria en la tabla de la base de datos
	Name        string  `json:"name" gorm:"size:60" validate:"required,gte=2,lte=60"`     // Nombre del producto, obligatorio, mínimo 2 caracteres, máximo 60 caracteres
	Description string  `json:"description,omitempty" gorm:"size:250" validate:"lte=250"` // Descripción "larga" del producto, no obligatorio
	Unit        string  `json:"unit" gorm:"size=32" validate:"required"`                  // Unidad de medida del producto (unidad, metros, litros, etc), hasta 32 caracteres, obligatorio
	Price       float64 `json:"price" validate:"required"`                                // Precio, obligatorio
	Stock       float64 `json:"stock,omitempty"`                                          // Cantidad del producto en stock
	IsActive    bool    `json:"is_active"`                                                // Indica si el producto está activo. Sólo para utilizar algún atributo de tipo boolean ;-)
}

// Repository representa el repositorio permanente de los productos.
// Se utiliza el concepto de interface para desacoplar la implementación específica del repositorio.
// Los métodos son los básicos de un ABM. Luego, en los usecases, quizás aparezcan otros métodos que se agregan y "extienden" esta interface.
// En general, los métodos devuelven los datos del producto afectado, en los casos de alta, consulta y actualización exitosos. Y un error en caso de falla.
type Repository interface {
	Create(product *Product) (*Product, error) // Create permite agregar un producto nuevo al repositorio
	GetByID(id uint) (*Product, error)         // GetByID permite recuperar un único producto, si existe, del repositorio
	GetByName(name string) (*Product, error)   // GetByName permite recuperar un único producto por nombre
	GetAll() ([]*Product, error)               // GetAll permite recuperar, en un slice, todos los productos existentes en el repositorio
	Update(product *Product) (*Product, error) // Update permite actualizar los datos de un producto
	Delete(product *Product) error             // Delete elmimina un producto del repositorio
}
