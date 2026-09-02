package catalog

import (
	"context"

	"github.com/mytheresa/go-hiring-challenge/domain/catalog"
)

func (app *App) GetProduct(ctx context.Context, code catalog.ProductCode) (*catalog.Product, error) {
	return app.productsRepo.GetProductByCode(code)
}
