package models

import (
	"github.com/mytheresa/go-hiring-challenge/domain/catalog"
	"github.com/shopspring/decimal"
)

// product represents a product in the catalog.
// It includes a unique code and a price.
type product struct {
	ID         uint            `gorm:"primaryKey"`
	Code       string          `gorm:"uniqueIndex;not null"`
	Price      decimal.Decimal `gorm:"type:decimal(10,2);not null"`
	Variants   []variant       `gorm:"foreignKey:ProductID"`
	CategoryID *uint
	Category   *category
}

func (p *product) TableName() string {
	return "products"
}

func productToDomain(dbProduct product) *catalog.Product {
	return catalog.RestoreProduct(dbProduct.Code, dbProduct.Price,
		catalog.ProductWithCategory(categoryToDomain(dbProduct.Category)),
		catalog.ProductWithVariants(variantsToDomain(dbProduct.Variants)...),
	)
}

func productsToDomain(dbProducts []product) []*catalog.Product {
	products := make([]*catalog.Product, len(dbProducts))
	for i, dbProduct := range dbProducts {
		products[i] = productToDomain(dbProduct)
	}
	return products
}
