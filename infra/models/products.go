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

func productToDomain(dbProduct product) (*catalog.Product, error) {
	id, err := catalog.ParseProductCode(dbProduct.Code)
	if err != nil {
		return nil, err
	}

	cat, err := categoryToDomain(dbProduct.Category)
	if err != nil {
		return nil, err
	}

	return catalog.RestoreProduct(id, dbProduct.Price,
		catalog.ProductWithCategory(cat),
		catalog.ProductWithVariants(variantsToDomain(dbProduct.Variants)...),
	), nil
}

func productsToDomain(dbProducts []product) ([]*catalog.Product, error) {
	products := make([]*catalog.Product, len(dbProducts))
	for i, dbProduct := range dbProducts {
		p, err := productToDomain(dbProduct)
		if err != nil {
			return nil, err
		}
		products[i] = p
	}
	return products, nil
}
