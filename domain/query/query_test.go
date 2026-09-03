package query_test

import (
	"testing"

	"github.com/mytheresa/go-hiring-challenge/domain/query"
	"github.com/stretchr/testify/assert"
)

func TestQuery(t *testing.T) {
	q := query.New().
		AddFilter("field1", query.Eq, "value1").
		AddFilter("field2", query.Ne, 10).
		AddFilter("field3", query.Gt, 20).
		AddFilter("field4", query.Gte, 30).
		AddFilter("field5", query.Lt, 40).
		AddFilter("field6", query.Lte, 50).
		AddSort("name", false).
		AddSort("email", true).
		AddPagination(1, 10)

	t.Run("filters", func(t *testing.T) {
		filter := q.Filters()[0]
		assert.Equal(t, "field1", filter.Field())
		assert.Equal(t, query.Eq, filter.Operator())
		assert.Equal(t, "value1", filter.Value())

		filter = q.Filters()[1]
		assert.Equal(t, "field2", filter.Field())
		assert.Equal(t, query.Ne, filter.Operator())
		assert.Equal(t, 10, filter.Value())

		filter = q.Filters()[2]
		assert.Equal(t, "field3", filter.Field())
		assert.Equal(t, query.Gt, filter.Operator())
		assert.Equal(t, 20, filter.Value())

		filter = q.Filters()[3]
		assert.Equal(t, "field4", filter.Field())
		assert.Equal(t, query.Gte, filter.Operator())
		assert.Equal(t, 30, filter.Value())

		filter = q.Filters()[4]
		assert.Equal(t, "field5", filter.Field())
		assert.Equal(t, query.Lt, filter.Operator())
		assert.Equal(t, 40, filter.Value())

		filter = q.Filters()[5]
		assert.Equal(t, "field6", filter.Field())
		assert.Equal(t, query.Lte, filter.Operator())
		assert.Equal(t, 50, filter.Value())
	})

	t.Run("sorting", func(t *testing.T) {
		sort := q.Sorts()[0]
		assert.Equal(t, "name", sort.Field())
		assert.False(t, sort.Desc())

		sort = q.Sorts()[1]
		assert.Equal(t, "email", sort.Field())
		assert.True(t, sort.Desc())
	})

	t.Run("pagination", func(t *testing.T) {
		pagination := q.Pagination()
		assert.Equal(t, 1, pagination.Page())
		assert.Equal(t, 10, pagination.PageSize())
	})

	t.Run("has filter", func(t *testing.T) {
		hasFilter := q.HasFilter("field1")
		assert.True(t, hasFilter)

		hasFilter = q.HasFilter("notExistentFilter")
		assert.False(t, hasFilter)
	})

	t.Run("get filter", func(t *testing.T) {
		filter := q.GetFilter("field1")
		assert.NotNil(t, filter)
		assert.Equal(t, "field1", filter.Field())
		assert.Equal(t, query.Eq, filter.Operator())
		assert.Equal(t, "value1", filter.Value())

		filter = q.GetFilter("non_existent_filter")
		assert.Nil(t, filter)
	})

	t.Run("add multiplefilters", func(t *testing.T) {

		q := query.New().AddFilters(
			query.NewFilter("field1", query.Eq, "value1"),
			query.NewFilter("field2", query.Ne, 10),
		)

		filter := q.Filters()[0]
		assert.Equal(t, "field1", filter.Field())
		assert.Equal(t, query.Eq, filter.Operator())
		assert.Equal(t, "value1", filter.Value())

		filter = q.Filters()[1]
		assert.Equal(t, "field2", filter.Field())
		assert.Equal(t, query.Ne, filter.Operator())
		assert.Equal(t, 10, filter.Value())
	})

}
