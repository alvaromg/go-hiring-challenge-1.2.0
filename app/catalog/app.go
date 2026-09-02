package catalog

type App struct {
	productsRepo   ProductsRepository
	categoriesRepo CategoriesRepository
}

func NewCatalogApp(productsRepo ProductsRepository, categoriesRepo CategoriesRepository) *App {
	return &App{
		productsRepo:   productsRepo,
		categoriesRepo: categoriesRepo,
	}
}
