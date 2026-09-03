package models

import (
	"context"
	"errors"
	"fmt"

	"github.com/mytheresa/go-hiring-challenge/domain/catalog"
	domainerrors "github.com/mytheresa/go-hiring-challenge/domain/errors"
	"github.com/mytheresa/go-hiring-challenge/domain/list"
	"github.com/mytheresa/go-hiring-challenge/domain/query"
	libmodel "github.com/mytheresa/go-hiring-challenge/infra/database"
	"github.com/mytheresa/go-hiring-challenge/infra/field"
	"gorm.io/gorm"
)

type ProductsRepository struct {
	db            *gorm.DB
	fieldsMapping map[string]string
}

func NewProductsRepository(db *gorm.DB) *ProductsRepository {
	return &ProductsRepository{
		db: db,
		fieldsMapping: map[string]string{
			field.Category: "c.code",
		},
	}
}
func (r *ProductsRepository) GetProducts(ctx context.Context, q *query.Query) (list.ListResponse[*catalog.Product], error) {
	var zero list.ListResponse[*catalog.Product]

	db := r.db.WithContext(ctx)

	// strict query validation
	validator := query.NewValidator().
		AllowFilter(field.Category, []query.Operator{query.Eq}, query.ValidateString).
		AllowFilter(field.Price, []query.Operator{query.Lt}, query.ValidatePrice).
		AllowSort(field.Code, field.Price)
	if err := validator.Validate(q); err != nil {
		return zero, err
	}

	// ensure the category filter, if present, refers to an existing category
	if q.HasFilter(field.Category) {
		categoryCode, ok := q.GetFilter(field.Category).Value().(string)
		if !ok {
			return zero, fmt.Errorf("%w: error parsing category field", libmodel.ErrorPersistence)
		}

		var count int64
		if err := db.Model(&category{}).Where("code = ?", categoryCode).Count(&count).Error; err != nil {
			return zero, fmt.Errorf("%w: error checking category existence: %s", libmodel.ErrorPersistence, err)
		}
		if count == 0 {
			return zero, fmt.Errorf("%w: category %q not found", domainerrors.ErrorNotFound, categoryCode)
		}
	}

	// build query to return just products ids
	db = db.Distinct("products.id", "products.code", "products.price")
	if q.HasFilter(field.Category) {
		db = db.Joins("JOIN categories c ON c.id = products.category_id")
	}

	// apply only filters for counting
	db = libmodel.ApplyQueryFilters(db, q, r.fieldsMapping)

	// count rows using only filters
	var total int64
	err := db.Model(&product{}).Count(&total).Error
	if err != nil {
		return zero, fmt.Errorf("%w: error counting products: %s", libmodel.ErrorPersistence, err)
	}
	if total == 0 {
		return zero, nil
	}

	// apply sort and pagination in addition to filters to find products in page
	db = libmodel.ApplyQuerySorts(db, q, r.fieldsMapping)
	db = libmodel.ApplyQueryPagination(db, q)

	// get products ids for requested page (filters + sort + pagination)
	// note: the select list includes products.code (required by Postgres so
	// DISTINCT queries can ORDER BY it), so we must scan into a struct that
	// has both columns rather than a plain []uint.
	var rows []product
	if err = db.Model(&product{}).Find(&rows).Error; err != nil {
		return zero, fmt.Errorf("%w: error finding products: %s", libmodel.ErrorPersistence, err)
	}
	ids := make([]uint, len(rows))
	for i, row := range rows {
		ids[i] = row.ID
	}

	// get full products for requested page, based on previous filtered, sortd and paginated ids
	var products []product
	db = r.db.WithContext(ctx)
	db = libmodel.ApplyQuerySorts(db, q, r.fieldsMapping)

	if err := db.WithContext(ctx).Preload("Variants").Preload("Category").Where("products.id IN ?", ids).Find(&products).Error; err != nil {
		return zero, fmt.Errorf("%w: error retrieving products: %s", libmodel.ErrorPersistence, err)
	}

	// build page response with domain products
	domainProducts, err := productsToDomain(products)
	if err != nil {
		return zero, fmt.Errorf("%w: error mapping products: %s", libmodel.ErrorPersistence, err)
	}
	listResp := list.NewListResponse[*catalog.Product]()
	listResp.SetTotal(uint(total))
	listResp.AddItems(domainProducts...)

	return listResp, nil
}

func (r *ProductsRepository) GetProductByCode(ctx context.Context, code catalog.ProductCode) (*catalog.Product, error) {
	var p product
	err := r.db.WithContext(ctx).Preload("Variants").Preload("Category").Where("code = ?", code).First(&p).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: product %q not found", domainerrors.ErrorNotFound, code)
		}
		return nil, fmt.Errorf("%w: error retrieving product: %s", libmodel.ErrorPersistence, err)
	}

	domainProduct, err := productToDomain(p)
	if err != nil {
		return nil, fmt.Errorf("%w: error mapping product: %s", libmodel.ErrorPersistence, err)
	}

	return domainProduct, nil
}
