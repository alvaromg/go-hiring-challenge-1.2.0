package catalog

import (
	"fmt"
	"unicode/utf8"

	domainerrors "github.com/mytheresa/go-hiring-challenge/domain/errors"
)

const maxCategoryNameLength = 256

type Category struct {
	code string
	name string
}

func (c *Category) Code() string {
	return c.code
}

func (c *Category) Name() string {
	return c.name
}

// NewCategory creates a new Category, enforcing its domain invariants.
func NewCategory(code, name string) (*Category, error) {
	if code == "" {
		return nil, fmt.Errorf("%w: category code must not be empty", domainerrors.ErrorDomainValidation)
	}
	if name == "" {
		return nil, fmt.Errorf("%w: category name must not be empty", domainerrors.ErrorDomainValidation)
	}
	if utf8.RuneCountInString(name) > maxCategoryNameLength {
		return nil, fmt.Errorf("%w: category name must not be longer than %d characters", domainerrors.ErrorDomainValidation, maxCategoryNameLength)
	}

	return &Category{
		code: code,
		name: name,
	}, nil
}

// RestoreCategory reconstructs a Category from already validated data, e.g. when loading it from persistence.
func RestoreCategory(code, name string) *Category {
	category := &Category{
		code: code,
		name: name,
	}

	return category
}
