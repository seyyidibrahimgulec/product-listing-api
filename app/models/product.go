package models

import (
	"github.com/seyyidibrahimgulec/product-listing/app/config"
	"gorm.io/gorm"
)

var db *gorm.DB

type Product struct {
	gorm.Model
	Name            string            `json:"name"`
	Price           float64           `json:"price"`
	Description     string            `json:"description"`
	Currency        string            `json:"currency"`
	SellerID        string            `json:"seller_id"`
	InStock         bool              `json:"in_stock"`
	DeliveryOptions []*DeliveryOption `gorm:"many2many:product_delivery_options;"json:"delivery_options"`
}

type DeliveryOption struct {
	gorm.Model
	Name     string     `json:"name"`
	Price    float64    `json:"price"`
	Currency string     `json:"currency"`
	Products []*Product `gorm:"many2many:product_delivery_options;"json:"-"`
}

func init() {
	db = config.ConnectDB()
	db.AutoMigrate(&Product{})
	db.AutoMigrate(&DeliveryOption{})
}

func (d *DeliveryOption) getOrCreateDeliveryOption() (*DeliveryOption, error) {
	var deliveryOption DeliveryOption
	if err := db.Where("Name=? AND Price=? AND Currency=?", d.Name, d.Price, d.Currency).First(&deliveryOption).Error; err != nil {
		if create_err := db.Create(&d).Error; create_err != nil {
			return nil, create_err
		}
		return d, nil
	}
	return &deliveryOption, nil
}

func (p *Product) appendDeliveryOptions(delivery_options []*DeliveryOption) {
	for _, do := range delivery_options {
		delivery_option, _ := do.getOrCreateDeliveryOption()
		db.Model(&p).Association("DeliveryOptions").Append(delivery_option)
	}
}

func GetAllProducts(offset, limit int, order_by string) []Product {
	var products []Product
	db.Preload("DeliveryOptions").Order(order_by).Offset(offset).Limit(limit).Find(&products)
	return products
}

func (p *Product) CreateProduct() (*Product, error) {
	// FIXME: -------------------------------------
	// For prevent delivery options duplication
	// we must clear p.DeliveryOptions list.
	// This is not best practice but omit delivery options
	// are not working this is a temp solution for now
	delivery_options := p.DeliveryOptions
	p.DeliveryOptions = []*DeliveryOption{}

	if err := db.Omit("DeliveryOptions").Create(&p).Error; err != nil {
		return nil, err
	}

	p.appendDeliveryOptions(delivery_options)

	return p, nil
}

func GetProductById(id int64) (*Product, error) {
	var product Product
	if err := db.Preload("DeliveryOptions").Where("ID=?", id).First(&product).Error; err != nil {
		return nil, err
	}

	return &product, nil
}

func (p *Product) UpdateProduct() (*Product, error) {
	// FIXME: -------------------------------------
	// For prevent delivery options duplication
	// we must clear p.DeliveryOptions list.
	// This is not best practice but omit delivery options
	// are not working this is a temp solution for now
	delivery_options := p.DeliveryOptions
	p.DeliveryOptions = []*DeliveryOption{}

	if err := db.Omit("DeliveryOptions").Save(&p).Error; err != nil {
		return nil, err
	}

	if delivery_options != nil {
		db.Model(&p).Association("DeliveryOptions").Clear()
		p.appendDeliveryOptions(delivery_options)
	}

	return p, nil
}
