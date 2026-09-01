package models

import (
	"github.com/mytheresa/go-hiring-challenge/domain/catalog"
	"github.com/shopspring/decimal"
)

// Variant represents a product variant in the catalog.
// It includes a unique name, SKU, and an optional price.
// Variants can be used to represent different configurations or options for a product.
type Variant struct {
	ID        uint            `gorm:"primaryKey"`
	ProductID uint            `gorm:"not null"`
	Name      string          `gorm:"not null"`
	SKU       string          `gorm:"uniqueIndex;not null"`
	Price     decimal.Decimal `gorm:"type:decimal(10,2);null"`
}

func (v *Variant) TableName() string {
	return "product_variants"
}

func variantToDomain(dbVariant Variant) catalog.Variant {
	return catalog.Variant{
		ID:        dbVariant.ID,
		ProductID: dbVariant.ProductID,
		Name:      dbVariant.Name,
		SKU:       dbVariant.SKU,
		Price:     dbVariant.Price,
	}
}

func variantsToDomain(dbVariants []Variant) []catalog.Variant {
	variants := make([]catalog.Variant, len(dbVariants))
	for i, dbVariant := range dbVariants {
		variants[i] = variantToDomain(dbVariant)
	}
	return variants
}
