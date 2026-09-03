package catalog

import (
	"fmt"
	"strconv"
	"strings"

	domainerrors "github.com/mytheresa/go-hiring-challenge/domain/errors"
)

const (
	productCodePrefix   = "PROD"
	maxProductCodeValue = 1000000000
)

// ProductCode is a value object representing a product identifier, formatted as "PROD###".
type ProductCode struct {
	value uint
}

func (cc ProductCode) Equal(other ProductCode) bool {
	return cc.value == other.value
}

// NewProductCode creates a new ProductCode, enforcing its domain invariants.
func NewProductCode(value uint) (ProductCode, error) {
	if value > maxProductCodeValue {
		return ProductCode{}, fmt.Errorf("%w: product code must not be greater than %d", domainerrors.ErrorDomainValidation, maxProductCodeValue)
	}

	return ProductCode{value: value}, nil
}

// ParseProductCode parses a formatted product code, e.g. "PROD001", into a ProductCode.
func ParseProductCode(code string) (ProductCode, error) {
	if !strings.HasPrefix(code, productCodePrefix) {
		return ProductCode{}, fmt.Errorf("%w: product code %q must start with %q", domainerrors.ErrorDomainValidation, code, productCodePrefix)
	}

	value, err := strconv.ParseUint(strings.TrimPrefix(code, productCodePrefix), 10, 64)
	if err != nil {
		return ProductCode{}, fmt.Errorf("%w: product code %q must end with a number", domainerrors.ErrorDomainValidation, code)
	}

	return NewProductCode(uint(value))
}

// String returns the formatted product code, e.g. "PROD001".
func (code ProductCode) String() string {
	return fmt.Sprintf("%s%03d", productCodePrefix, code.value)
}
