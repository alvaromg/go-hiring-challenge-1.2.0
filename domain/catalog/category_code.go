package catalog

import (
	"fmt"
	"strconv"
	"strings"

	domainerrors "github.com/mytheresa/go-hiring-challenge/domain/errors"
)

const (
	categoryCodePrefix   = "CAT"
	maxCategoryCodeValue = 999
)

// CategoryCode is a value object representing a category identifier, formatted as "CAT###".
type CategoryCode struct {
	value uint
}

// NewCategoryCode creates a new CategoryCode, enforcing its domain invariants.
func NewCategoryCode(value uint) (CategoryCode, error) {
	if value > maxCategoryCodeValue {
		return CategoryCode{}, fmt.Errorf("%w: category code must not be greater than %d", domainerrors.ErrorDomainValidation, maxCategoryCodeValue)
	}

	return CategoryCode{value: value}, nil
}

// ParseCategoryCode parses a formatted category code, e.g. "CAT001", into a CategoryCode.
func ParseCategoryCode(code string) (CategoryCode, error) {
	if !strings.HasPrefix(code, categoryCodePrefix) {
		return CategoryCode{}, fmt.Errorf("%w: category code %q must start with %q", domainerrors.ErrorDomainValidation, code, categoryCodePrefix)
	}

	value, err := strconv.ParseUint(strings.TrimPrefix(code, categoryCodePrefix), 10, 64)
	if err != nil {
		return CategoryCode{}, fmt.Errorf("%w: category code %q must end with a number", domainerrors.ErrorDomainValidation, code)
	}

	return NewCategoryCode(uint(value))
}

// Value returns the underlying numeric identifier.
func (code CategoryCode) Value() uint {
	return code.value
}

// String returns the formatted category code, e.g. "CAT001".
func (code CategoryCode) String() string {
	return fmt.Sprintf("%s%03d", categoryCodePrefix, code.value)
}
