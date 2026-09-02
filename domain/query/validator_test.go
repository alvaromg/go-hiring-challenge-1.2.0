package query_test

import (
	"testing"

	"github.com/mytheresa/go-hiring-challenge/domain/query"
	"github.com/stretchr/testify/assert"
)

func TestQueryValidator(t *testing.T) {
	// Create a validator with rules
	validator := query.NewValidator().
		AllowFilter("name", query.StringOperators(), query.ValidateString).
		AllowFilter("age", query.NumericOperators(), query.ValidateInt).
		AllowFilter("status", query.JustEqualOperator(), query.ValidateString).
		AllowSort("name", "age")

	t.Run("valid query", func(t *testing.T) {
		query := query.New().
			AddFilter("name", query.Eq, "John").
			AddFilter("age", query.Gt, 18).
			AddFilter("status", query.Eq, "active").
			AddSort("name", false).
			AddSort("age", true)

		err := validator.Validate(query)
		assert.NoError(t, err)
	})

	t.Run("invalid field", func(t *testing.T) {
		query := query.New().
			AddFilter("invalid_field", query.Eq, "value")

		err := validator.Validate(query)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "filtering by field \"invalid_field\" is not allowed")
	})

	t.Run("invalid operator", func(t *testing.T) {
		query := query.New().
			AddFilter("name", query.Gt, "John")

		err := validator.Validate(query)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "operator \"gt\" is not allowed for field \"name\"")
	})

	t.Run("invalid value type", func(t *testing.T) {
		query := query.New().
			AddFilter("age", query.Eq, "not a number")

		err := validator.Validate(query)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid value for field \"age\"")
	})

	t.Run("invalid sort field", func(t *testing.T) {
		query := query.New().
			AddSort("invalid_sort", false)

		err := validator.Validate(query)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "sorting by field \"invalid_sort\" is not allowed")
	})

	t.Run("nil query", func(t *testing.T) {
		err := validator.Validate(nil)
		assert.NoError(t, err)
	})
}

func TestHelperFunctions(t *testing.T) {

	t.Run("NumericOperators", func(t *testing.T) {
		ops := query.NumericOperators()
		assert.Contains(t, ops, query.Eq)
		assert.Contains(t, ops, query.Ne)
		assert.Contains(t, ops, query.Gt)
		assert.Contains(t, ops, query.Gte)
		assert.Contains(t, ops, query.Lt)
		assert.Contains(t, ops, query.Lte)
	})

	t.Run("StringOperators", func(t *testing.T) {
		ops := query.StringOperators()
		assert.Contains(t, ops, query.Eq)
		assert.Contains(t, ops, query.Ne)
		assert.NotContains(t, ops, query.Gt)
	})

	t.Run("ValidateString", func(t *testing.T) {
		assert.NoError(t, query.ValidateString("test"))
		assert.Error(t, query.ValidateString(123))
	})

	t.Run("ValidateInt", func(t *testing.T) {
		assert.NoError(t, query.ValidateInt(123))
		assert.Error(t, query.ValidateInt("test"))
		assert.Error(t, query.ValidateInt(12.3))
	})

	t.Run("ValidatePrice", func(t *testing.T) {
		assert.NoError(t, query.ValidatePrice("19.99"))
		assert.Error(t, query.ValidatePrice("not a price"))
		assert.Error(t, query.ValidatePrice(19.99))
	})

	t.Run("ValidateBool", func(t *testing.T) {
		assert.NoError(t, query.ValidateBool(true))
		assert.NoError(t, query.ValidateBool(false))
		assert.NoError(t, query.ValidateBool("true"))
		assert.NoError(t, query.ValidateBool("false"))
		assert.NoError(t, query.ValidateBool("1"))
		assert.NoError(t, query.ValidateBool("0"))
		assert.Error(t, query.ValidateBool("not a bool"))
		assert.Error(t, query.ValidateBool(123))
	})

	t.Run("ValidateNumeric", func(t *testing.T) {
		assert.NoError(t, query.ValidateNumeric(123))
		assert.NoError(t, query.ValidateNumeric(12.3))
		assert.Error(t, query.ValidateNumeric("test"))
	})

}
