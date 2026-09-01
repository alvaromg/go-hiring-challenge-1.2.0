package query_test

import (
	"testing"

	"github.com/mytheresa/go-hiring-challenge/domain/query"
	"github.com/stretchr/testify/assert"
)

// Test enum for validation testing
type testStatus struct{}

func (testStatus) ValidValues() []string {
	return []string{"active", "inactive", "pending"}
}

func TestQueryValidator(t *testing.T) {
	// Create a validator with rules
	validator := query.NewValidator().
		AllowFilter("name", query.StringOperators(), query.ValidateString).
		AllowFilter("age", query.NumericOperators(), query.ValidateInt).
		RequireFilter("tenant_id", query.EqualityOperators(), query.ValidateString).
		AllowSort("name", "age")

	t.Run("valid query", func(t *testing.T) {
		query := query.New().
			AddFilter("name", query.Eq, "John").
			AddFilter("age", query.Gt, 18).
			AddFilter("status", query.Eq, "active").
			AddFilter("tenant_id", query.Eq, "tenant123").
			AddSort("name", false).
			AddSort("age", true)

		err := validator.Validate(query)
		assert.NoError(t, err)
	})

	t.Run("invalid field", func(t *testing.T) {
		query := query.New().
			AddFilter("invalid_field", query.Eq, "value").
			AddFilter("tenant_id", query.Eq, "tenant123")

		err := validator.Validate(query)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "filtering by field \"invalid_field\" is not allowed")
	})

	t.Run("invalid operator", func(t *testing.T) {
		query := query.New().
			AddFilter("name", query.Gt, "John").
			AddFilter("tenant_id", query.Eq, "tenant123")

		err := validator.Validate(query)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "operator \"gt\" is not allowed for field \"name\"")
	})

	t.Run("invalid value type", func(t *testing.T) {
		query := query.New().
			AddFilter("age", query.Eq, "not a number").
			AddFilter("tenant_id", query.Eq, "tenant123")

		err := validator.Validate(query)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid value for field \"age\"")
	})

	t.Run("invalid enum value", func(t *testing.T) {
		query := query.New().
			AddFilter("status", query.Eq, "invalid_status").
			AddFilter("tenant_id", query.Eq, "tenant123")

		err := validator.Validate(query)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid value for field \"status\"")
	})

	t.Run("invalid sort field", func(t *testing.T) {
		query := query.New().
			AddFilter("tenant_id", query.Eq, "tenant123").
			AddSort("invalid_sort", false)

		err := validator.Validate(query)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "sorting by field \"invalid_sort\" is not allowed")
	})

	t.Run("nil query", func(t *testing.T) {
		err := validator.Validate(nil)
		assert.NoError(t, err)
	})

	t.Run("empty query", func(t *testing.T) {
		query := query.New()
		err := validator.Validate(query)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "missing required filter for field \"tenant_id\"")
	})

	t.Run("missing required filter", func(t *testing.T) {
		query := query.New().
			AddFilter("name", query.Eq, "John").
			AddFilter("age", query.Gt, 18).
			AddSort("name", false)

		err := validator.Validate(query)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "missing required filter for field \"tenant_id\"")
	})

	t.Run("valid In operator", func(t *testing.T) {
		query := query.New().
			AddFilter("status", query.In, []string{"active"}).
			AddFilter("tenant_id", query.Eq, "tenant123")

		err := validator.Validate(query)
		assert.NoError(t, err)
	})

	t.Run("valid Nin operator with slice validator", func(t *testing.T) {
		// Create a validator that allows slice operations
		sliceValidator := query.NewValidator().
			AllowFilter("ids", query.SliceOperators(), query.ValidateSlice[string]).
			RequireFilter("tenant_id", query.EqualityOperators(), query.ValidateString)

		query := query.New().
			AddFilter("ids", query.Nin, []string{"id1", "id2"}).
			AddFilter("tenant_id", query.Eq, "tenant123")

		err := sliceValidator.Validate(query)
		assert.NoError(t, err)
	})
}

func TestHelperFunctions(t *testing.T) {
	t.Run("AllOperatorsExcept", func(t *testing.T) {
		ops := query.AllOperatorsExcept(query.In, query.Nin, query.Like)
		assert.NotContains(t, ops, query.In)
		assert.NotContains(t, ops, query.Nin)
		assert.NotContains(t, ops, query.Like)
		assert.Contains(t, ops, query.Eq)
		assert.Contains(t, ops, query.Gt)
	})

	t.Run("NumericOperators", func(t *testing.T) {
		ops := query.NumericOperators()
		assert.Contains(t, ops, query.Eq)
		assert.Contains(t, ops, query.Ne)
		assert.Contains(t, ops, query.Gt)
		assert.Contains(t, ops, query.Gte)
		assert.Contains(t, ops, query.Lt)
		assert.Contains(t, ops, query.Lte)
		assert.NotContains(t, ops, query.Like)
	})

	t.Run("StringOperators", func(t *testing.T) {
		ops := query.StringOperators()
		assert.Contains(t, ops, query.Eq)
		assert.Contains(t, ops, query.Ne)
		assert.Contains(t, ops, query.Like)
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

	t.Run("ValidateNumeric", func(t *testing.T) {
		assert.NoError(t, query.ValidateNumeric(123))
		assert.NoError(t, query.ValidateNumeric(12.3))
		assert.Error(t, query.ValidateNumeric("test"))
	})

	t.Run("SliceOperators", func(t *testing.T) {
		ops := query.SliceOperators()
		assert.Contains(t, ops, query.In)
		assert.Contains(t, ops, query.Nin)
		assert.NotContains(t, ops, query.Eq)
		assert.NotContains(t, ops, query.Like)
		assert.Len(t, ops, 2)
	})

	t.Run("EqualityOperators includes Nin", func(t *testing.T) {
		ops := query.EqualityOperators()
		assert.Contains(t, ops, query.Eq)
		assert.Contains(t, ops, query.Ne)
		assert.Contains(t, ops, query.In)
		assert.Contains(t, ops, query.Nin)
		assert.NotContains(t, ops, query.Like)
		assert.NotContains(t, ops, query.Gt)
	})

	t.Run("ValidateAllowedValues", func(t *testing.T) {
		validator := query.ValidateAllowedValues([]any{"active", "inactive"})
		assert.NoError(t, validator("active"))
		assert.Error(t, validator("pending"))
		assert.Error(t, validator(123))
	})

}
