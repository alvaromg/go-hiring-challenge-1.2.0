package models

import (
	"github.com/mytheresa/go-hiring-challenge/domain/catalog"
)

// category represents a product category in the catalog.
type category struct {
	ID   uint   `gorm:"primaryKey"`
	Code string `gorm:"uniqueIndex;not null"`
	Name string `gorm:"not null"`
}

func (c *category) TableName() string {
	return "categories"
}

func categoryToDomain(dbCategory *category) (*catalog.Category, error) {
	if dbCategory == nil {
		return nil, nil
	}

	id, err := catalog.ParseCategoryCode(dbCategory.Code)
	if err != nil {
		return nil, err
	}

	return catalog.RestoreCategory(id, dbCategory.Name), nil
}

func categoriesToDomain(dbCategories []category) ([]*catalog.Category, error) {
	categories := make([]*catalog.Category, len(dbCategories))
	for i, dbCategory := range dbCategories {
		c, err := categoryToDomain(&dbCategory)
		if err != nil {
			return nil, err
		}
		categories[i] = c
	}
	return categories, nil
}

func categoryFromDomain(c *catalog.Category) category {
	return category{
		Code: c.Code().String(),
		Name: c.Name(),
	}
}

func categoriesFromDomain(categories []*catalog.Category) []category {
	dbCategories := make([]category, len(categories))
	for i, c := range categories {
		dbCategories[i] = categoryFromDomain(c)
	}
	return dbCategories
}
