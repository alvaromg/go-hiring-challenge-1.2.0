package catalog

import "github.com/mytheresa/go-hiring-challenge/shared"

type App struct {
	productsRepo   ProductsRepository
	categoriesRepo CategoriesRepository
	logger         shared.Logger
}

func NewCatalogApp(logger shared.Logger, productsRepo ProductsRepository, categoriesRepo CategoriesRepository) *App {
	return &App{
		productsRepo:   productsRepo,
		categoriesRepo: categoriesRepo,
		logger:         logger,
	}
}
