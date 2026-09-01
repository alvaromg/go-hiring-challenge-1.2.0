package api

import (
	"net/http"

	"github.com/mytheresa/go-hiring-challenge/domain/catalog"
	"github.com/shopspring/decimal"
)

type product struct {
	ID    uint            `json:"id"`
	Code  string          `json:"code"`
	Price decimal.Decimal `json:"price"`
	// Variants []Variant
}

type productsList struct {
	Products []product `json:"products"`
}

func decodeGetProductsListRequest(r *http.Request) (struct{}, error) {
	return struct{}{}, nil
}

func encodeProductResponse(p *catalog.Product) product {
	return product{
		ID:    p.ID,
		Code:  p.Code,
		Price: p.Price,
		// Variants: p.Variants,
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
