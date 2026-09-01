package models

import (
	"github.com/mytheresa/go-hiring-challenge/domain/catalog"
	"github.com/shopspring/decimal"
)

// Product represents a product in the catalog.
// It includes a unique code and a price.
type Product struct {
	ID         uint            `gorm:"primaryKey"`
	Code       string          `gorm:"uniqueIndex;not null"`
	Price      decimal.Decimal `gorm:"type:decimal(10,2);not null"`
	Variants   []Variant       `gorm:"foreignKey:ProductID"`
	Categories []Category      `gorm:"many2many:product_categories;"`
}

func (p *Product) TableName() string {
	return "products"
}

func productToDomain(dbProduct Product) *catalog.Product {
	product := &catalog.Product{
		Code:       dbProduct.Code,
		Price:      dbProduct.Price,
		Variants:   variantsToDomain(dbProduct.Variants),
		Categories: categoriesToDomain(dbProduct.Categories),
	}
	return product
}

func productsToDomain(dbProducts []Product) []*catalog.Product {
	products := make([]*catalog.Product, len(dbProducts))
	for i, dbProduct := range dbProducts {
		products[i] = productToDomain(dbProduct)
	}
	return products
}
