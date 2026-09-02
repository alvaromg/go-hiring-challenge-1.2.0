package models

import (
	"fmt"

	"github.com/mytheresa/go-hiring-challenge/domain/catalog"
	domainerrors "github.com/mytheresa/go-hiring-challenge/domain/errors"
	"github.com/mytheresa/go-hiring-challenge/domain/list"
	"github.com/mytheresa/go-hiring-challenge/domain/query"
	libmodel "github.com/mytheresa/go-hiring-challenge/lib/model"
	"gorm.io/gorm"
)

const (
	FieldPrice    = "price"
	FieldCategory = "category"
)

type ProductsRepository struct {
	db            *gorm.DB
	fieldsMapping map[string]string
}

func NewProductsRepository(db *gorm.DB) *ProductsRepository {
	return &ProductsRepository{
		db: db,
		fieldsMapping: map[string]string{
			FieldCategory: "c.code",
		},
	}
}
func (r *ProductsRepository) GetAllProducts(q *query.Query) (list.ListResponse[*catalog.Product], error) {
	var listRes list.ListResponse[*catalog.Product]

	// strict query validation
	validator := query.NewValidator().
		AllowFilter(FieldCategory, []query.Operator{query.Eq}, query.ValidateString).
		AllowFilter(FieldPrice, []query.Operator{query.Lt}, query.ValidateString)
	if err := validator.Validate(q); err != nil {
		return listRes, err
	}

	// ensure the category filter, if present, refers to an existing category
	if q.HasFilter(FieldCategory) {
		categoryCode, ok := q.GetFilter(FieldCategory).Value().(string)
		if !ok {
			return listRes, fmt.Errorf("%w: error parsing category field", libmodel.ErrorPersistence)
		}

		var count int64
		if err := r.db.Model(&category{}).Where("code = ?", categoryCode).Count(&count).Error; err != nil {
			return listRes, fmt.Errorf("%w: error checking category existence: %s", libmodel.ErrorPersistence, err)
		}
		if count == 0 {
			return listRes, fmt.Errorf("%w: category %q not found", domainerrors.ErrorNotFound, categoryCode)
		}
	}

	// build query to return just ids
	db := r.db.Distinct("products.id")
	if q.HasFilter(FieldCategory) {
		db = db.Joins("JOIN categories c ON c.id = products.category_id")
	}

	// apply only filters for counting
	db = libmodel.ApplyQueryFilters(db, q, r.fieldsMapping)

	// count rows using only filters
	var total int64
	err := db.Model(&product{}).Count(&total).Error
	if err != nil {
		return listRes, fmt.Errorf("%w: error counting products: %s", libmodel.ErrorPersistence, err)
	}
	if total == 0 {
		return listRes, nil
	}

	// apply sort and pagination in addition to filters to find products in page
	db = libmodel.ApplyQuerySorts(db, q, r.fieldsMapping)
	db = libmodel.ApplyQueryPagination(db, q)

	// get products ids for requested page (filters + sort + pagination)
	var ids []uint
	if err = db.Model(&product{}).Find(&ids).Error; err != nil {
		return listRes, fmt.Errorf("%w: error finding products: %s", libmodel.ErrorPersistence, err)
	}

	// get full products for requested page, based on previous filtered, sortd and paginated ids
	var products []product
	if err := r.db.Preload("Variants").Preload("Category").Where("products.id IN ?", ids).Find(&products).Error; err != nil {
		return listRes, fmt.Errorf("%w: error retrieving products: %s", libmodel.ErrorPersistence, err)
	}

	// build page response with domain products
	listRes.SetTotal(uint(total))
	listRes.AddItems(productsToDomain(products)...)

	return listRes, nil
}
