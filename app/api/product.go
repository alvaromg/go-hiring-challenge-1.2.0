package api

import (
	"fmt"
	"net/http"

	"github.com/mytheresa/go-hiring-challenge/domain/catalog"
	"github.com/mytheresa/go-hiring-challenge/lib/rest"
	"github.com/shopspring/decimal"
)

type variant struct {
	Name  string           `json:"name"`
	SKU   string           `json:"sku"`
	Price *decimal.Decimal `json:"price"`
}

type product struct {
	Code     string          `json:"code"`
	Price    decimal.Decimal `json:"price"`
	Category *category       `json:"category"`
	Variants []variant       `json:"variants"`
}

type productsData struct {
	Products []product `json:"products"`
}

func wrapProductsResponse(products []product) productsData {
	return productsData{Products: products}
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
	if p.Category() != nil {
		c := encodeCategoryResponse(p.Category())
		cat = &c
	}

	variants := make([]variant, len(p.Variants()))
	for i, v := range p.Variants() {
		variants[i] = encodeVariantResponse(v)
	}

	return product{
		Code:     p.Code(),
		Price:    p.Price(),
		Category: cat,
		Variants: variants,
	}
}

func encodeProductDataResponse(p *catalog.Product) (product, error) {
	return encodeProductResponse(p), nil
}

func decodeProductCodeFromRequest(r *http.Request) (string, error) {
	code := r.PathValue("code")
	if code == "" {
		return "", fmt.Errorf("%w: missing product code", rest.ErrorBadRequest)
	}
	return code, nil
}
