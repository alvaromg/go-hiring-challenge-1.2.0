package catalog

import (
	"context"

	"github.com/mytheresa/go-hiring-challenge/domain/catalog"
)

// RequestedCategory carries the raw data for a single category to be created.
type RequestedCategory struct {
	Code catalog.CategoryCode
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

	created, err := app.categoriesRepo.CreateCategories(ctx, categories)
	if err != nil {
		return nil, err
	}

	codes := make([]string, len(created))
	for i, c := range created {
		codes[i] = c.Code().String()
	}
	app.logger.WithContext(ctx).WithField("codes", codes).Infof("categories created")

	return created, nil
}
