package api

import (
	"fmt"
	"net/http"

	"github.com/mytheresa/go-hiring-challenge/app/catalog"
	domainCatalog "github.com/mytheresa/go-hiring-challenge/domain/catalog"
	"github.com/mytheresa/go-hiring-challenge/infra/rest"
	"github.com/mytheresa/go-hiring-challenge/shared"
)

// ProductsList
// @Summary      Product detail
// @Tags         Catalog
// @Router       /v1/catalog/{code} [get]
// @Accept       json
// @Produce      json
// @Param        code   path      string  true  "Product code"
// @Success      200  {object}  rest.Response[product]
// @Failure      500  {object}  rest.ErrorResponse
// @Failure      400  {object}  rest.ErrorResponse
// @Failure      404  {object}  rest.ErrorResponse
func productsDetailController(monitor shared.Monitor, app *catalog.App) http.HandlerFunc {
	return rest.NewHandler(monitor, app.ProductDetail, decodeProductCodeFromRequest, encodeProductDetailResponse, http.StatusOK)
}

type variant struct {
	Name  string `json:"name"`
	SKU   string `json:"sku"`
	Price string `json:"price"`
}

type product struct {
	Code     string    `json:"code"`
	Price    string    `json:"price"`
	Category *category `json:"category"`
	Variants []variant `json:"variants"`
}

func encodeVariantResponse(v *domainCatalog.Variant) variant {
	return variant{
		Name:  v.Name(),
		SKU:   v.SKU(),
		Price: v.Price().String(),
	}
}

func encodeProductResponse(p *domainCatalog.Product) product {
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
		Price:    p.Price().String(),
		Category: cat,
		Variants: variants,
	}
}

func decodeProductCodeFromRequest(r *http.Request) (domainCatalog.ProductCode, error) {
	var zero domainCatalog.ProductCode

	codeStr := r.PathValue("code")
	if codeStr == "" {
		return zero, fmt.Errorf("%w: missing product code", rest.ErrorBadRequest)
	}

	productCode, err := domainCatalog.ParseProductCode(codeStr)
	if err != nil {
		return zero, fmt.Errorf("%w: %s", rest.ErrorBadRequest, err)
	}

	return productCode, nil
}

func encodeProductDetailResponse(p *domainCatalog.Product) (product, error) {
	return encodeProductResponse(p), nil
}
