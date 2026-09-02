package models

import (
	"fmt"

	"github.com/mytheresa/go-hiring-challenge/domain/catalog"
	domainerrors "github.com/mytheresa/go-hiring-challenge/domain/errors"
	"github.com/mytheresa/go-hiring-challenge/domain/list"
	"github.com/mytheresa/go-hiring-challenge/domain/query"
	libmodel "github.com/mytheresa/go-hiring-challenge/infra/database"
	"gorm.io/gorm"
)

type CategoriesRepository struct {
	db *gorm.DB
}

func NewCategoriesRepository(db *gorm.DB) *CategoriesRepository {
	return &CategoriesRepository{
		db: db,
	}
}

func (r *CategoriesRepository) GetCategories(q *query.Query) (list.ListResponse[*catalog.Category], error) {
	var listRes list.ListResponse[*catalog.Category]

	// strict query validation: only pagination is allowed, no filters or sorts
	validator := query.NewValidator()
	if err := validator.Validate(q); err != nil {
		return listRes, err
	}

	// count all categories
	var total int64
	if err := r.db.Model(&category{}).Count(&total).Error; err != nil {
		return listRes, fmt.Errorf("%w: error counting categories: %s", libmodel.ErrorPersistence, err)
	}
	if total == 0 {
		return listRes, nil
	}

	// apply pagination to find categories in page
	db := libmodel.ApplyQueryPagination(r.db, q)

	var categories []category
	if err := db.Find(&categories).Error; err != nil {
		return listRes, fmt.Errorf("%w: error retrieving categories: %s", libmodel.ErrorPersistence, err)
	}

	// build page response with domain categories
	domainCategories, err := categoriesToDomain(categories)
	if err != nil {
		return listRes, fmt.Errorf("%w: error mapping categories: %s", libmodel.ErrorPersistence, err)
	}
	listRes.SetTotal(uint(total))
	listRes.AddItems(domainCategories...)

	return listRes, nil
}

// CreateCategories persists all given categories in a single transaction.
func (r *CategoriesRepository) CreateCategories(categories []*catalog.Category) ([]*catalog.Category, error) {
	dbCategories := categoriesFromDomain(categories)

	err := r.db.Transaction(func(tx *gorm.DB) error {
		return tx.Create(&dbCategories).Error
	})
	if err != nil {
		if value, ok := libmodel.AsDuplicateKeyError(err); ok {
			return nil, fmt.Errorf("%w: category with code %q already exists", domainerrors.ErrorDuplicatedResource, value)
		}
		return nil, fmt.Errorf("%w: error creating categories: %s", libmodel.ErrorPersistence, err)
	}

	domainCategories, err := categoriesToDomain(dbCategories)
	if err != nil {
		return nil, fmt.Errorf("%w: error mapping categories: %s", libmodel.ErrorPersistence, err)
	}

	return domainCategories, nil
}
