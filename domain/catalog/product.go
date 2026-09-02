package catalog

import "github.com/shopspring/decimal"

type Product struct {
	code     string
	price    decimal.Decimal
	variants []Variant
	category *Category
}

func (p *Product) Code() string {
	return p.code
}

func (p *Product) Price() decimal.Decimal {
	return p.price
}

func (p *Product) Category() *Category {
	return p.category
}

func (p *Product) Variants() []Variant {
	return p.variants
}

// fillVariantsPrices sets product's price to all variants with price = nil
func (p *Product) fillVariantsPrices() {
	productPrice := p.Price()
	for i := range p.variants {
		if p.variants[i].Price == nil {

			p.variants[i].Price = &productPrice
		}
	}
}

// ProductOption implements the functiona option pattern for products
type ProductOption func(*Product)

func RestoreProduct(code string, price decimal.Decimal, opts ...ProductOption) *Product {
	product := &Product{
		code:     code,
		price:    price,
		variants: []Variant{},
		category: nil,
	}

	for _, opt := range opts {
		opt(product)
	}

	product.fillVariantsPrices()
	return product
}

func ProductWithCategory(c *Category) ProductOption {
	return func(p *Product) {
		p.category = c
	}
}

func ProductWithVariants(v ...Variant) ProductOption {
	return func(p *Product) {
		p.variants = append(p.variants, v...)
	}
}
