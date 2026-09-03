package catalog

import "github.com/mytheresa/go-hiring-challenge/domain/price"

type Variant struct {
	name  string
	sku   string
	price *price.Price
}

func (v *Variant) Name() string {
	return v.name
}

func (v *Variant) SKU() string {
	return v.sku
}

func (v *Variant) Price() *price.Price {
	return v.price
}

func (v *Variant) ChangePrice(p *price.Price) {
	v.price = p
}

func RestoreVariant(name, sku string, price *price.Price) *Variant {
	return &Variant{
		name:  name,
		sku:   sku,
		price: price,
	}
}
