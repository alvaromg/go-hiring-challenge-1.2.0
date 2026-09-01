package model

import (
	"fmt"
	"strings"

	"github.com/mytheresa/go-hiring-challenge/domain/query"
	"gorm.io/gorm"
)

// ApplyQuery applies the query criteria to a GORM DB instance
func ApplyQuery(db *gorm.DB, q *query.Query, fieldsMapping map[string]string) *gorm.DB {
	if q == nil {
		return db
	}

	db = ApplyQueryFilters(db, q, fieldsMapping)
	db = ApplyQuerySorts(db, q, fieldsMapping)
	db = ApplyQueryPagination(db, q)

	return db
}

func mappedField(fieldsMapping map[string]string, field string) string {
	if mappedField, ok := fieldsMapping[field]; ok {
		return mappedField
	}
	return field
}

// quoteIdent quotes a (possibly qualified, e.g. "c.code") identifier so each
// part is quoted individually, producing "c"."code" rather than "c.code".
func quoteIdent(field string) string {
	parts := strings.Split(field, ".")
	for i, p := range parts {
		parts[i] = fmt.Sprintf("%q", p)
	}
	return strings.Join(parts, ".")
}

func ApplyQueryFilters(db *gorm.DB, q *query.Query, fieldsMapping map[string]string) *gorm.DB {
	if q == nil {
		return db
	}

	// Apply filters
	for _, filter := range q.Filters() {
		field := quoteIdent(mappedField(fieldsMapping, filter.Field()))
		switch filter.Operator() {
		case query.Eq:
			db = db.Where(fmt.Sprintf("%s = ?", field), filter.Value())
		case query.Ne:
			db = db.Where(fmt.Sprintf("%s != ?", field), filter.Value())
		case query.Gt:
			db = db.Where(fmt.Sprintf("%s > ?", field), filter.Value())
		case query.Gte:
			db = db.Where(fmt.Sprintf("%s >= ?", field), filter.Value())
		case query.Lt:
			db = db.Where(fmt.Sprintf("%s < ?", field), filter.Value())
		case query.Lte:
			db = db.Where(fmt.Sprintf("%s <= ?", field), filter.Value())
		case query.In:
			db = db.Where(fmt.Sprintf("%s IN (?)", field), filter.Value())
		case query.Nin:
			db = db.Where(fmt.Sprintf("%s NOT IN (?)", field), filter.Value())
		case query.Like:
			db = db.Where(fmt.Sprintf("%s ILIKE ?", field), fmt.Sprintf("%%%s%%", filter.Value()))
		case query.Is:
			if filter.Value() == nil {
				db = db.Where(fmt.Sprintf("%s IS NULL", field))
			} else {
				db = db.Where(fmt.Sprintf("%s IS ?", field), filter.Value())
			}
		case query.IsNot:
			if filter.Value() == nil {
				db = db.Where(fmt.Sprintf("%s IS NOT NULL", field))
			} else {
				db = db.Where(fmt.Sprintf("%s IS NOT ?", field), filter.Value())
			}
		}
	}

	return db
}

func ApplyQuerySorts(db *gorm.DB, q *query.Query, fieldsMapping map[string]string) *gorm.DB {
	if q == nil {
		return db
	}

	for _, sort := range q.Sorts() {
		field := quoteIdent(mappedField(fieldsMapping, sort.Field()))
		if sort.Desc() {
			db = db.Order(fmt.Sprintf("%s DESC", field))
		} else {
			db = db.Order(field)
		}
	}

	return db
}

func ApplyQueryPagination(db *gorm.DB, q *query.Query) *gorm.DB {
	if q == nil {
		return db
	}

	if q.Pagination() != nil {
		offset := (q.Pagination().Page() - 1) * q.Pagination().PageSize()
		db = db.Offset(offset).Limit(q.Pagination().PageSize())
	}
	return db
}
