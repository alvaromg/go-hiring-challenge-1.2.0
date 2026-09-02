package query_test

import (
	"testing"

	"github.com/mytheresa/go-hiring-challenge/domain/query"
	"github.com/stretchr/testify/assert"
)

func TestNewQuery(t *testing.T) {
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

	sort := q.Sorts()[0]
	assert.Equal(t, "name", sort.Field())
	assert.False(t, sort.Desc())

	sort = q.Sorts()[1]
	assert.Equal(t, "email", sort.Field())
	assert.True(t, sort.Desc())

	pagination := q.Pagination()
	assert.Equal(t, 1, pagination.Page())
	assert.Equal(t, 10, pagination.PageSize())
}

func TestOperatorParsing(t *testing.T) {
	t.Run("parse valid operators", func(t *testing.T) {
		testCases := []struct {
			input    string
			expected query.Operator
		}{
			{"eq", query.Eq},
			{"ne", query.Ne},
			{"gt", query.Gt},
			{"gte", query.Gte},
			{"lt", query.Lt},
			{"lte", query.Lte},
		}

		for _, tc := range testCases {
			t.Run(tc.input, func(t *testing.T) {
				op, err := query.ParseOperator(tc.input)
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, op)
			})
		}
	})

	t.Run("parse invalid operator", func(t *testing.T) {
		_, err := query.ParseOperator("invalid")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid operator invalid")
	})

	t.Run("operator validation", func(t *testing.T) {
		assert.False(t, query.Operator("invalid").IsValid())
	})
}
