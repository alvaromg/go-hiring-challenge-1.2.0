package query_test

import (
	"testing"

	"github.com/mytheresa/go-hiring-challenge/domain/query"
	"github.com/stretchr/testify/assert"
)

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
		assert.Contains(t, err.Error(), "invalid query: invalid operator \"invalid\"")
	})

	t.Run("operator validation", func(t *testing.T) {
		assert.False(t, query.Operator("invalid").IsValid())
	})
}
