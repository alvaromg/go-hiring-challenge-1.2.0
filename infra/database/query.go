package database

import (
	"fmt"
	"strings"

	"github.com/mytheresa/go-hiring-challenge/domain/query"
	"gorm.io/gorm"
)

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
