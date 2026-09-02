package query

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"time"

	"github.com/mytheresa/go-hiring-challenge/domain/price"
)

var ErrorInvalidQuery = errors.New("invalid query")

// filterValidation defines validation rules for a specific filter field
type filterValidation struct {
	// field is the name of the field being validated
	field string
	// allowedOperators is the list of allowed operators for this field
	allowedOperators []Operator
	// validateValue is a function that validates the filter value
	// Return nil if the value is valid, or an error describing why it's invalid
	validateValue func(value any) error
	// required indicates if this filter must be present in the query
	required bool
}

// sortValidation defines validation rules for a specific sort field
type sortValidation struct {
	// field is the name of the field being sorted
	field string
}

// validator provides validation for Query objects
type validator struct {
	// allowedFilters defines which filters are allowed and their validation rules
	allowedFilters []filterValidation
	// allowedSorts defines which sort fields are allowed
	allowedSorts []sortValidation
}

// const (
// 	ErrTypeInvalidQuery liberr.ErrType = "ErrTypeInvalidQuery"
// )

func newErrInvalidQuery(format string, args ...any) error {
	return fmt.Errorf("%w: %s", ErrorInvalidQuery, fmt.Sprintf(format, args...))
}

// NewValidator creates a new QueryValidator
func NewValidator() *validator {
	return &validator{
		allowedFilters: make([]filterValidation, 0),
		allowedSorts:   make([]sortValidation, 0),
	}
}

// AllowFilter adds a filter validation rule
func (v *validator) AllowFilter(field string, allowedOperators []Operator, validateValue func(value any) error) *validator {
	v.allowedFilters = append(v.allowedFilters, filterValidation{
		field:            field,
		allowedOperators: allowedOperators,
		validateValue:    validateValue,
		required:         false,
	})
	return v
}

// RequireFilter adds a mandatory filter validation rule
func (v *validator) RequireFilter(field string, allowedOperators []Operator, validateValue func(value any) error) *validator {
	v.allowedFilters = append(v.allowedFilters, filterValidation{
		field:            field,
		allowedOperators: allowedOperators,
		validateValue:    validateValue,
		required:         true,
	})
	return v
}

// AllowSort adds a sort validation rule
func (v *validator) AllowSort(field ...string) *validator {
	for _, f := range field {
		v.allowedSorts = append(v.allowedSorts, sortValidation{
			field: f,
		})
	}
	return v
}

// Validate validates a query against the defined rules
func (v *validator) Validate(query *Query) error {
	if query == nil {
		return nil
	}

	// Check for required filters
	err := v.validateRequiredFilters(query)
	if err != nil {
		return err
	}

	// Validate filters
	for _, filter := range query.Filters() {
		if err := v.validateFilter(filter); err != nil {
			return err
		}
	}

	// Validate sorts
	for _, sort := range query.Sorts() {
		if err := v.validateSort(sort); err != nil {
			return err
		}
	}

	return nil
}

// validateRequiredFilters checks if all required filters are present
func (v *validator) validateRequiredFilters(query *Query) error {
	// Build a map of provided filter fields for quick lookup
	providedFields := make(map[string]bool)
	for _, filter := range query.Filters() {
		providedFields[filter.Field()] = true
	}

	// Check if required filters are present
	for _, filterValidation := range v.allowedFilters {
		if filterValidation.required && !providedFields[filterValidation.field] {
			return newErrInvalidQuery("missing required filter for field %q", filterValidation.field)
		}
	}

	return nil
}

// validateFilter validates a single filter against the defined rules
func (v *validator) validateFilter(filter Filter) error {
	// Check if the field is allowed
	filterValidation, err := v.findFilterValidation(filter.Field())
	if err != nil {
		return err
	}

	// Check if the operator is allowed for this field
	if !slices.Contains(filterValidation.allowedOperators, filter.Operator()) {
		return newErrInvalidQuery("operator %q is not allowed for field %q", filter.Operator(), filter.Field())
	}

	// Validate the value
	if filterValidation.validateValue != nil {
		if err := filterValidation.validateValue(filter.Value()); err != nil {
			return newErrInvalidQuery("invalid value for field %q: %v", filter.Field(), err)
		}
	}

	return nil
}

// validateSort validates a single sort against the defined rules
func (v *validator) validateSort(sort Sort) error {
	// Check if the sort field is allowed
	for _, allowedSort := range v.allowedSorts {
		if allowedSort.field == sort.Field() {
			return nil
		}
	}
	return newErrInvalidQuery("sorting by field %q is not allowed", sort.Field())
}

// findFilterValidation finds the validation rules for a given field
func (v *validator) findFilterValidation(field string) (filterValidation, error) {
	for _, filterValidation := range v.allowedFilters {
		if filterValidation.field == field {
			return filterValidation, nil
		}
	}
	return filterValidation{}, newErrInvalidQuery("filtering by field %q is not allowed", field)
}

// Helper functions for common value validations

// ValidateString validates that the value is a string
func ValidateString(value any) error {
	if _, ok := value.(string); !ok {
		return errors.New("value must be a string")
	}
	return nil
}

func ValidateDate(value any) error {
	dateStr, ok := value.(string)
	if !ok {
		return newErrInvalidQuery("value must be a date string")
	}

	if _, err := time.Parse(time.DateOnly, dateStr); err != nil {
		return newErrInvalidQuery("invalid date format %q", dateStr)
	}
	return nil
}

// ValidatePrice validates that the value is a valid price string, e.g. "19.99"
func ValidatePrice(value any) error {
	str, ok := value.(string)
	if !ok {
		return errors.New("value must be a valid price")
	}

	if _, err := price.Parse(str); err != nil {
		return fmt.Errorf("price %q is not valid", value)
	}

	return nil
}

// ValidateBool validates that the value is a boolean or valid boolean string
func ValidateBool(value any) error {
	// Check if it's already a boolean
	if _, ok := value.(bool); ok {
		return nil
	}

	// Check if it's a valid boolean string
	if str, ok := value.(string); ok {
		_, err := strconv.ParseBool(str)
		if err != nil {
			return errors.New("value must be a boolean or valid boolean string (true, false, 1, 0)")
		}
		return nil
	}

	return errors.New("value must be a boolean or valid boolean string")
}

// ValidateInt validates that the value is an integer
func ValidateInt(value any) error {
	switch value.(type) {
	case int, int8, int16, int32, int64:
		return nil
	default:
		return errors.New("value must be an integer")
	}
}

// ValidateFloat validates that the value is a float
func ValidateFloat(value any) error {
	switch value.(type) {
	case float32, float64:
		return nil
	default:
		return errors.New("value must be a float")
	}
}

// ValidateNumeric validates that the value is a number (integer or float)
func ValidateNumeric(value any) error {
	switch value.(type) {
	case int, int8, int16, int32, int64, float32, float64:
		return nil
	default:
		return errors.New("value must be a number")
	}
}

// NumericOperators returns operators suitable for numeric values
func NumericOperators() []Operator {
	return []Operator{Eq, Ne, Gt, Gte, Lt, Lte}
}

// StringOperators returns operators suitable for string values
func StringOperators() []Operator {
	return []Operator{Eq, Ne}
}

// EqualityOperators returns operators for equality comparisons

func JustEqualOperator() []Operator {
	return []Operator{Eq}
}
