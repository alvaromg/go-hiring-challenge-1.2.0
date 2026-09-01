package rest

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/mytheresa/go-hiring-challenge/domain/query"
)

const (
	pageParam       = "page"
	pageSizeParam   = "pageSize"
	defaultPageSize = 10
	maxPageSize     = 100
)

func DecodeQueryFromRequest(r *http.Request) (*query.Query, error) {
	// Extract query parameters
	params := r.URL.Query()

	q := query.New()
	var err error

	// Pagination handling
	pageStr := params.Get(pageParam)
	page := 1
	if pageStr != "" {
		page, err = strconv.Atoi(pageStr)
		if err != nil {
			return nil, fmt.Errorf("invalid page parameter %q", pageStr)
		}
		if page <= 0 {
			return nil, fmt.Errorf("page paremeter must be greater or equal than 0")
		}
	}

	pageSize := defaultPageSize
	if params.Get(pageSizeParam) != "" {
		pageSize, err = strconv.Atoi(params.Get(pageSizeParam))
		if err != nil {
			return nil, fmt.Errorf("invalid page size parameter %q", pageStr)
		}
		if pageSize < 1 {
			return nil, fmt.Errorf("page size paremeter must be greater or equal than 1")
		}
		if pageSize > maxPageSize {
			return nil, fmt.Errorf("page size paremeter must lower or equal than %d", maxPageSize)
		}
	}
	q = q.AddPagination(page, pageSize)

	// Sorting
	sortParam := params.Get("sort")
	if sortParam != "" {
		sorts := strings.Split(sortParam, ",")
		for _, sort := range sorts {
			desc := strings.HasPrefix(sort, "-")
			if desc {
				sort = strings.TrimPrefix(sort, "-")
			}
			q = q.AddSort(sort, desc)
		}
	}

	filters, err := parseFilters(r.URL.Query())
	if err != nil {
		return nil, fmt.Errorf("unable to parse page filters: %s", err)
	}
	q = q.AddFilters(filters...)

	return q, nil
}

const (
	filterField = "filter"
)

func parseFilters(values url.Values) ([]query.Filter, error) {
	parsedFilters := make(map[string]map[string]string)
	filters := []query.Filter{}

	for key, val := range values {
		if !strings.HasPrefix(key, fmt.Sprintf("%s_", filterField)) {
			continue // Ignore non-filter parameters
		}

		// Extract field and operator
		parts := strings.Split(key, "_")
		if len(parts) != 3 {
			continue // Skip invalid keys
		}

		field, operator := parts[1], parts[2]
		if parsedFilters[field] == nil {
			parsedFilters[field] = make(map[string]string)
		}

		op, err := query.ParseOperator(operator)
		if err != nil {
			return nil, fmt.Errorf("unable to parse query operator %q", operator)
		}

		// Parse filter values based on operator type (preserve original operator)
		var filterValue any

		switch op {
		case query.In, query.Nin:
			// Parse comma-separated values into a slice
			values := strings.Split(val[0], ",")
			// Trim whitespace from each value
			for i := range values {
				values[i] = strings.TrimSpace(values[i])
			}
			filterValue = values
		case query.Eq, query.Ne:
			// For Eq/Ne operators, use single string value
			filterValue = strings.TrimSpace(val[0])
		default:
			// For other operators (Like, Gt, Lt, etc.), use single value
			filterValue = strings.TrimSpace(val[0])
		}

		parsedFilters[field][operator] = val[0] // Store original value for reference
		filters = append(filters, query.NewFilter(field, op, filterValue))
	}

	return filters, nil
}
