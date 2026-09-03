package catalog_test

type categoryDTO struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type variantDTO struct {
	Name  string  `json:"name"`
	SKU   string  `json:"sku"`
	Price *string `json:"price"`
}

type productDTO struct {
	Code     string       `json:"code"`
	Price    string       `json:"price"`
	Category *categoryDTO `json:"category"`
	Variants []variantDTO `json:"variants"`
}

type productDetailDTO struct {
	Product productDTO `json:"product"`
}

type productsListDTO struct {
	Products []productDTO `json:"products"`
}

type categoriesListDTO struct {
	Categories []categoryDTO `json:"categories"`
}
