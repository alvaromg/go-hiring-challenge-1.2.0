package models

import (
	"github.com/mytheresa/go-hiring-challenge/domain/catalog"
	"github.com/mytheresa/go-hiring-challenge/domain/price"
	"github.com/shopspring/decimal"
)

// variant represents a product variant in the catalog.
// It includes a unique name, SKU, and an optional price.
// Variants can be used to represent different configurations or options for a product.
type variant struct {
	ID        uint             `gorm:"primaryKey"`
	ProductID uint             `gorm:"not null"`
	Name      string           `gorm:"not null"`
	SKU       string           `gorm:"uniqueIndex;not null"`
	Price     *decimal.Decimal `gorm:"type:decimal(10,2);null"`
}

func (v *variant) TableName() string {
	return "product_variants"
}

func variantToDomain(dbVariant variant) *catalog.Variant {
	var variantPrice *price.Price
	if dbVariant.Price != nil {
		p := price.New(*dbVariant.Price)
		variantPrice = &p
	}

	return catalog.RestoreVariant(dbVariant.Name, dbVariant.SKU, variantPrice)
}

func variantsToDomain(dbVariants []variant) []*catalog.Variant {
	variants := make([]*catalog.Variant, len(dbVariants))
	for i, dbVariant := range dbVariants {
		variants[i] = variantToDomain(dbVariant)
	}
	return variants
}
