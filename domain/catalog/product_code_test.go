package catalog_test

import (
	"testing"

	"github.com/mytheresa/go-hiring-challenge/domain/catalog"
	domainerrors "github.com/mytheresa/go-hiring-challenge/domain/errors"
	"github.com/stretchr/testify/assert"
)

func TestNewProductCode(t *testing.T) {
	t.Run("new valid code", func(t *testing.T) {
		code, err := catalog.NewProductCode(1)
		assert.NoError(t, err)
		assert.Equal(t, "PROD001", code.String())
	})

	t.Run("code value greater than max", func(t *testing.T) {
		_, err := catalog.NewProductCode(10000000001)
		assert.Error(t, err)
		assert.ErrorIs(t, err, domainerrors.ErrorDomainValidation)
	})
}

func TestParseProductCode(t *testing.T) {
	t.Run("parse valid code", func(t *testing.T) {
		codeStr := "PROD123"
		code, err := catalog.ParseProductCode(codeStr)
		assert.NoError(t, err)
		assert.Equal(t, codeStr, code.String())
	})

	t.Run("parse invalid code", func(t *testing.T) {
		codeStr := "invalid"
		_, err := catalog.ParseProductCode(codeStr)
		assert.Error(t, err)
		assert.ErrorIs(t, err, domainerrors.ErrorDomainValidation)
	})

	t.Run("parse invalid code value", func(t *testing.T) {
		codeStr := "PRODinvalid"
		_, err := catalog.ParseProductCode(codeStr)
		assert.Error(t, err)
		assert.ErrorIs(t, err, domainerrors.ErrorDomainValidation)
	})

}
