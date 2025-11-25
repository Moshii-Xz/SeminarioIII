package catalog

import (
	"errors"
	"time"

	"github.com/mordmora/expirapp/internal/domain"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(product *domain.Product) error {
	return r.db.Create(product).Error
}

func (r *Repository) FindByID(id uint) (*domain.Product, error) {
	var product domain.Product
	err := r.db.Preload("Category").Preload("Store").First(&product, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}

	return &product, nil
}

func (r *Repository) Update(product *domain.Product) error {
	// Usar Select para especificar explícitamente qué campos actualizar
	// Esto asegura que los campos puntero se actualicen correctamente
	return r.db.Model(product).Select(
		"nombre",
		"descripcion",
		"imagen_url",
		"precio",
		"precio_original",
		"precio_descuento",
		"fecha_vencimiento",
		"stock",
		"etiqueta",
		"id_categoria",
	).Updates(product).Error
}

func (r *Repository) Delete(id uint) error {
	return r.db.Delete(&domain.Product{}, id).Error
}

// ProductWithStoreName es una estructura auxiliar para cargar el nombre de la tienda
type ProductWithStoreName struct {
	domain.Product
	StoreName string `gorm:"column:nombre_tienda"`
}

func (r *Repository) List(limit, offset int) ([]domain.Product, int64, error) {
	var products []domain.Product
	var total int64

	if err := r.db.Model(&domain.Product{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Hacemos un JOIN con la tabla usuario para obtener el nombre de la tienda
	// Asumimos que id_tienda en producto es FK a id_tienda en tienda, y tienda.id_tienda es FK a usuario.id_usuario
	// O más simple: id_tienda en producto apunta directamente al usuario que es la tienda (según el modelo actual StoreID es uint)

	// Consulta optimizada para traer el nombre de la tienda
	// Usamos Preload para Category, pero para StoreName necesitamos un Join manual o Preload inteligente si el modelo lo soportara
	// Dado que no queremos cambiar el modelo domain.Product para agregar el campo StoreName (no es persistente),
	// vamos a iterar y cargar los nombres o usar un mapa.

	// Opción eficiente: Cargar productos y luego cargar los nombres de las tiendas en batch
	err := r.db.Preload("Category").Limit(limit).Offset(offset).Find(&products).Error
	if err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

// GetStoreNameByID obtiene el nombre de la tienda (usuario) dado su ID
func (r *Repository) GetStoreNameByID(storeID uint) (string, error) {
	var name string
	// La tabla es 'usuario', el id es 'id_usuario'
	err := r.db.Table("usuario").Select("nombre").Where("id_usuario = ?", storeID).Scan(&name).Error
	return name, err
}

func (r *Repository) FindByName(name string) (*domain.Product, error) {
	var product domain.Product

	err := r.db.Where("nombre = ?", name).First(&product).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}

	return &product, nil
}

func (r *Repository) FindByExpirationDate(date time.Time) ([]domain.Product, error) {
	var products []domain.Product

	err := r.db.Where("fecha_vencimiento = ?", date).Find(&products).Error
	return products, err
}

func (r *Repository) FindExpiringSoon(days int) ([]domain.Product, error) {
	var products []domain.Product
	threshold := time.Now().AddDate(0, 0, days)

	err := r.db.Where("fecha_vencimiento <= ? AND fecha_vencimiento >= ?", threshold, time.Now()).Find(&products).Error
	return products, err
}

func (r *Repository) UpdateStock(id uint, quantity int) error {
	return r.db.Model(&domain.Product{}).Where("id_producto = ?", id).Update("stock", gorm.Expr("stock + ?", quantity)).Error
}
