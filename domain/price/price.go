// Package price provides a domain facade around the third-party decimal
// implementation used to represent monetary amounts, so that dependency
// does not need to be imported anywhere else in the codebase.
package price

import (
	"fmt"

	domainerrors "github.com/mytheresa/go-hiring-challenge/domain/errors"
	"github.com/shopspring/decimal"
)

// Price is a value object representing a monetary amount.
type Price struct {
	value decimal.Decimal
}

// New creates a Price from a decimal value.
func New(value decimal.Decimal) Price {
	return Price{value: value}
}

// Parse parses a decimal string, e.g. "19.99", into a Price.
func Parse(value string) (Price, error) {
	d, err := decimal.NewFromString(value)
	if err != nil {
		return Price{}, fmt.Errorf("%w: price %q must be a decimal number", domainerrors.ErrorDomainValidation, value)
	}

	return Price{value: d}, nil
}

// Decimal returns the underlying decimal value.
func (p Price) Decimal() decimal.Decimal {
	return p.value
}

// String returns the price formatted as a plain decimal string, e.g. "19.99".
func (p Price) String() string {
	return p.value.String()
}

// MarshalJSON implements json.Marshaler.
func (p Price) MarshalJSON() ([]byte, error) {
	return p.value.MarshalJSON()
}

// UnmarshalJSON implements json.Unmarshaler.
func (p *Price) UnmarshalJSON(data []byte) error {
	return p.value.UnmarshalJSON(data)
}
