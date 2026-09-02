package api

import (
	"fmt"
	"net/http"

	"github.com/mytheresa/go-hiring-challenge/domain/catalog"
	"github.com/mytheresa/go-hiring-challenge/domain/price"
	"github.com/mytheresa/go-hiring-challenge/infra/rest"
)

type variant struct {
	Name  string       `json:"name"`
	SKU   string       `json:"sku"`
	Price *price.Price `json:"price"`
}

type product struct {
	Code     string      `json:"code"`
	Price    price.Price `json:"price"`
	Category *category   `json:"category"`
	Variants []variant   `json:"variants"`
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
		Code:     p.Code().String(),
		Price:    p.Price(),
		Category: cat,
		Variants: variants,
	}
}

func decodeProductCodeFromRequest(r *http.Request) (catalog.ProductCode, error) {
	var zero catalog.ProductCode

	codeStr := r.PathValue("code")
	if codeStr == "" {
		return zero, fmt.Errorf("%w: missing product code", rest.ErrorBadRequest)
	}

	productCode, err := catalog.ParseProductCode(codeStr)
	if err != nil {
		return zero, fmt.Errorf("%w: %s", rest.ErrorBadRequest, err)
	}

	return productCode, nil
}

func encodeProductDetailResponse(p *catalog.Product) (product, error) {
	return encodeProductResponse(p), nil
}
