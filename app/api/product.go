package api

import (
	"net/http"

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

type productsList struct {
	Products []product `json:"products"`
}

func decodeGetProductsListRequest(r *http.Request) (struct{}, error) {
	return struct{}{}, nil
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

func encodeProductsListResponse(products []*catalog.Product) (productsList, error) {
	productsList := productsList{
		Products: make([]product, len(products)),
	}

	for i, p := range products {
		productsList.Products[i] = encodeProductResponse(p)
	}

	return productsList, nil
}
