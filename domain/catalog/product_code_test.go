package catalog_test

import (
	"testing"

	"github.com/mytheresa/go-hiring-challenge/domain/catalog"
	domainerrors "github.com/mytheresa/go-hiring-challenge/domain/errors"
	"github.com/stretchr/testify/assert"
)

func TestNewProductCode(t *testing.T) {
	code, err := catalog.NewProductCode(1)
	assert.NoError(t, err)
	assert.Equal(t, "PROD001", code.String())

	_, err = catalog.NewProductCode(1000)
	assert.Error(t, err)
	assert.ErrorIs(t, err, domainerrors.ErrorDomainValidation)
}

func TestParseProductCode(t *testing.T) {
	codeStr := "PROD123"
	code, err := catalog.ParseProductCode(codeStr)
	assert.NoError(t, err)
	assert.Equal(t, codeStr, code.String())

	codeStr = "invalid"
	_, err = catalog.ParseProductCode(codeStr)
	assert.Error(t, err)
	assert.ErrorIs(t, err, domainerrors.ErrorDomainValidation)

	codeStr = "PRODinvalid"
	_, err = catalog.ParseProductCode(codeStr)
	assert.Error(t, err)
	assert.ErrorIs(t, err, domainerrors.ErrorDomainValidation)

}
