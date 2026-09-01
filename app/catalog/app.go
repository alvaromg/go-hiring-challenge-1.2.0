package catalog

type App struct {
	productsRepo ProductsRepository
}

func NewCatalogApp(productsRepo ProductsRepository) *App {
	return &App{
		productsRepo: productsRepo,
	}
}
