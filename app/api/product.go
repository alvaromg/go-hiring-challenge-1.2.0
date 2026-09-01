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
	Code       string          `json:"code"`
	Price      decimal.Decimal `json:"price"`
	Categories []category      `json:"categories"`
	Variants   []variant       `json:"variants"`
}

func encodeVariantResponse(v catalog.Variant) variant {
	return variant{
		Name:  v.Name,
		SKU:   v.SKU,
		Price: v.Price,
	}
}

func encodeProductResponse(p *catalog.Product) product {
	categories := make([]category, len(p.Categories))
	for i, c := range p.Categories {
		categories[i] = encodeCategoryResponse(c)
	}

	variants := make([]variant, len(p.Variants))
	for i, v := range p.Variants {
		variants[i] = encodeVariantResponse(v)
	}

	return product{
		Code:       p.Code,
		Price:      p.Price,
		Categories: categories,
		Variants:   variants,
	}
}
