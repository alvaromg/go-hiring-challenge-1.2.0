package catalog_test

import (
	"testing"

	"github.com/mytheresa/go-hiring-challenge/domain/catalog"
	domainerrors "github.com/mytheresa/go-hiring-challenge/domain/errors"
	"github.com/stretchr/testify/assert"
)

func TestNewCategoryCode(t *testing.T) {
	code, err := catalog.NewCategoryCode(1)
	assert.NoError(t, err)
	assert.Equal(t, "CAT001", code.String())

	_, err = catalog.NewCategoryCode(1000)
	assert.Error(t, err)
	assert.ErrorIs(t, err, domainerrors.ErrorDomainValidation)
}

func TestParseCategoryCode(t *testing.T) {
	codeStr := "CAT123"
	code, err := catalog.ParseCategoryCode(codeStr)
	assert.NoError(t, err)
	assert.Equal(t, codeStr, code.String())

	codeStr = "invalid"
	_, err = catalog.ParseCategoryCode(codeStr)
	assert.Error(t, err)
	assert.ErrorIs(t, err, domainerrors.ErrorDomainValidation)

	codeStr = "CATinvalid"
	_, err = catalog.ParseCategoryCode(codeStr)
	assert.Error(t, err)
	assert.ErrorIs(t, err, domainerrors.ErrorDomainValidation)

}
