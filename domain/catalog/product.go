package catalog

import "github.com/mytheresa/go-hiring-challenge/domain/price"

type Product struct {
	code     ProductCode
	price    price.Price
	variants []*Variant
	category *Category
}

func (p *Product) Code() ProductCode {
	return p.code
}

func (p *Product) Price() price.Price {
	return p.price
}

func (p *Product) Category() *Category {
	return p.category
}

func (p *Product) Variants() []*Variant {
	return p.variants
}

func (p *Product) Equal(other *Product) bool {
	return p.code.Equal(other.Code())
}

// fillVariantsPrices sets product's price to all variants with price = nil
func (p *Product) fillVariantsPrices() {
	productPrice := p.Price()
	for i := range p.variants {
		if p.variants[i].Price() == nil {
			p.variants[i].ChangePrice(&productPrice)
		}
	}
}

// ProductOption implements the functional options pattern for products
type ProductOption func(*Product)

func RestoreProduct(code ProductCode, amount price.Price, opts ...ProductOption) *Product {
	product := &Product{
		code:     code,
		price:    amount,
		variants: []*Variant{},
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

func ProductWithVariants(v ...*Variant) ProductOption {
	return func(p *Product) {
		p.variants = append(p.variants, v...)
	}
}
