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

func categoryToDomain(dbCategory category) catalog.Category {
	return catalog.Category{
		Code: dbCategory.Code,
		Name: dbCategory.Name,
	}
}

func categoryToDomainPtr(dbCategory *category) *catalog.Category {
	if dbCategory == nil {
		return nil
	}
	c := categoryToDomain(*dbCategory)
	return &c
}
