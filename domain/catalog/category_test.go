package catalog_test

import (
	"testing"

	"github.com/mytheresa/go-hiring-challenge/domain/catalog"
	domainerrors "github.com/mytheresa/go-hiring-challenge/domain/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCategory(t *testing.T) {
	code, err := catalog.NewCategoryCode(1)
	require.NoError(t, err)

	t.Run("new valid category", func(t *testing.T) {
		cat, err := catalog.NewCategory(code, "Accesories")
		assert.NoError(t, err)
		assert.True(t, cat.Code().Equal(code))
		assert.Equal(t, "Accesories", cat.Name())
	})

	t.Run("new category with empty name", func(t *testing.T) {
		_, err := catalog.NewCategory(code, "")
		assert.Error(t, err)
		assert.ErrorIs(t, err, domainerrors.ErrorDomainValidation)
		assert.ErrorContains(t, err, "category name must not be empty")
	})

	t.Run("new category with name too long", func(t *testing.T) {
		_, err := catalog.NewCategory(code, "123456789012345678901234567890123")
		assert.Error(t, err)
		assert.ErrorIs(t, err, domainerrors.ErrorDomainValidation)
		assert.ErrorContains(t, err, "category name must not be longer than 32 characters")
	})

	t.Run("restore category", func(t *testing.T) {
		cat := catalog.RestoreCategory(code, "Accesories")
		assert.True(t, cat.Code().Equal(code))
		assert.Equal(t, "Accesories", cat.Name())
	})
}
