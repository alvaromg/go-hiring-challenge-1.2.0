package catalog

import (
	"context"

	"github.com/mytheresa/go-hiring-challenge/domain/catalog"
)

type App struct {
	productsRepo ProductsRepository
}

func NewCatalogApp(productsRepo ProductsRepository) *App {
	return &App{
		productsRepo: productsRepo,
	}
}

func (app *App) GetProducts(ctx context.Context, query struct{}) ([]*catalog.Product, error) {
	res, err := app.productsRepo.GetAllProducts()
	if err != nil {
		return nil, err
	}

	return res, nil
}
