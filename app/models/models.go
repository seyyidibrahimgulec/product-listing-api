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
	CurrencyID      int               `json:"currency_id"`
	Currency        CurrencyType      `json:"currency"`
	SellerID        int               `json:"seller_id"`
	Seller          Seller            `json:"seller"`
	InStock         bool              `json:"in_stock"`
	DeliveryOptions []DeliveryOptions `gorm:"many2many:product_delivery_options;"json:"delivery_options"`
}

type Seller struct {
	gorm.Model
	SellerID string `gorm:"unique"json:"seller_id"`
}

type DeliveryOptions struct {
	gorm.Model
	Name       string       `json:"name"`
	Price      float64      `json:"price"`
	CurrencyID int          `json:"currency_id"`
	Currency   CurrencyType `json:"currency"`
	Product    []Product    `gorm:"many2many:product_delivery_options;"json:"products"`
}

type CurrencyType struct {
	gorm.Model
	Code string `gorm:"unique"json:"code"`
}

func init() {
	db = config.ConnectDB()
	db.AutoMigrate(&Product{})
	db.AutoMigrate(&Seller{})
	db.AutoMigrate(&DeliveryOptions{})
	db.AutoMigrate(&CurrencyType{})
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
