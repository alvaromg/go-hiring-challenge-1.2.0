package catalog

import (
	"context"

	"github.com/mytheresa/go-hiring-challenge/domain/catalog"
)

// RequestedCategory carries the raw data for a single category to be created.
type RequestedCategory struct {
	Code string
	Name string
}

type CreateCategoriesInput struct {
	Categories []RequestedCategory
}

func (app *App) CreateCategories(ctx context.Context, input CreateCategoriesInput) ([]*catalog.Category, error) {
	categories := make([]*catalog.Category, len(input.Categories))
	for i, requested := range input.Categories {
		c, err := catalog.NewCategory(requested.Code, requested.Name)
		if err != nil {
			return nil, err
		}
		categories[i] = c
	}

	return app.categoriesRepo.CreateCategories(categories)
}
