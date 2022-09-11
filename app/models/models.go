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
	DeliveryOptions []DeliveryOptions `gorm:"many2many:product_delivery_options;"json:"delivery_options"`
}

type DeliveryOptions struct {
	gorm.Model
	Name     string    `json:"name"`
	Price    float64   `json:"price"`
	Currency string    `json:"currency"`
	Product  []Product `gorm:"many2many:product_delivery_options;"json:"products"`
}

func init() {
	db = config.ConnectDB()
	db.AutoMigrate(&Product{})
	db.AutoMigrate(&DeliveryOptions{})
}

func GetAllProducts() []Product {
	var Products []Product
	db.Find(&Products)
	return Products
}

func (p *Product) CreateProduct() (*Product, error) {
	if err := db.Create(&p).Error; err != nil {
		return nil, err
	}
	db.Create(&p)
	return p, nil
}

func GetProductById(id int64) (*Product, error) {
	var product Product
	if err := db.Where("ID=?", id).First(&product).Error; err != nil {
		return nil, err
	}

	return &product, nil
}

func (p *Product) UpdateProduct() (*Product, error) {
	if err := db.Save(&p).Error; err != nil {
		return nil, err
	}
	db.Save(&p)
	return p, nil
}
