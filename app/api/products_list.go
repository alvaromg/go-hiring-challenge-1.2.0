package api

type productsList struct {
	Products []product `json:"products"`
}

func wrapProductsListResponse(products []product) productsList {
	return productsList{Products: products}
}
