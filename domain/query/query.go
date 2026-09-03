package query

const (
	defaultPageSize = 20
)

// Query encapsulates all query parameters
type Query struct {
	filters    []Filter
	sorts      []Sort
	pagination *Pagination
}

// New creates a new Query instance
func New() *Query {
	defaultPagination := NewPagination(1, defaultPageSize)
	return &Query{
		filters:    make([]Filter, 0),
		sorts:      make([]Sort, 0),
		pagination: &defaultPagination,
	}
}

// AddFilter adds a new filter
func (q *Query) AddFilter(field string, op Operator, value any) *Query {
	q.filters = append(q.filters, NewFilter(field, op, value))
	return q
}

func (q *Query) HasFilter(field string) bool {
	for _, f := range q.filters {
		if f.Field() == field {
			return true
		}
	}
	return false
}

func (q *Query) HasAnySort() bool {
	return len(q.sorts) > 0
}

func (q *Query) GetFilter(field string) *Filter {
	for _, f := range q.filters {
		if f.Field() == field {
			return &f
		}
	}
	return nil
}

func (q *Query) AddFilters(filters ...Filter) *Query {
	q.filters = append(q.filters, filters...)
	return q
}

// AddSort adds a sort configuration
func (q *Query) AddSort(field string, desc bool) *Query {
	q.sorts = append(q.sorts, NewSort(field, desc))
	return q
}

// SetPagination sets pagination parameters
func (q *Query) AddPagination(page, pageSize int) *Query {
	*q.pagination = NewPagination(page, pageSize)
	return q
}

// Filters returns all criteria
func (q *Query) Filters() []Filter {
	return q.filters
}

// Sorts returns all sorts
func (q *Query) Sorts() []Sort {
	return q.sorts
}

// Pagination returns pagination configuration
func (q *Query) Pagination() *Pagination {
	return q.pagination
}
