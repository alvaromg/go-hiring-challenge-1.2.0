package api

import (
	"github.com/mytheresa/go-hiring-challenge/domain/catalog"
	"github.com/shopspring/decimal"
)

type variant struct {
	Name  string          `json:"name"`
	SKU   string          `json:"sku"`
	Price decimal.Decimal `json:"price"`
}

type product struct {
	Code     string          `json:"code"`
	Price    decimal.Decimal `json:"price"`
	Category *category       `json:"category"`
	Variants []variant       `json:"variants"`
}

func encodeVariantResponse(v catalog.Variant) variant {
	return variant{
		Name:  v.Name,
		SKU:   v.SKU,
		Price: v.Price,
	}
}

func encodeProductResponse(p *catalog.Product) product {
	var cat *category
	if p.Category != nil {
		c := encodeCategoryResponse(*p.Category)
		cat = &c
	}

	variants := make([]variant, len(p.Variants))
	for i, v := range p.Variants {
		variants[i] = encodeVariantResponse(v)
	}

	return product{
		Code:     p.Code,
		Price:    p.Price,
		Category: cat,
		Variants: variants,
	}
}
