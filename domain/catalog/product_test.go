package catalog_test

import (
	"testing"

	"github.com/mytheresa/go-hiring-challenge/domain/catalog"
	"github.com/mytheresa/go-hiring-challenge/domain/price"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProduct(t *testing.T) {
	prodCode, err := catalog.NewProductCode(1)
	require.NoError(t, err)
	prodPrice, err := price.Parse("12.99")
	require.NoError(t, err)

	t.Run("restore product", func(t *testing.T) {
		prod := catalog.RestoreProduct(prodCode, prodPrice)
		assert.Equal(t, prodCode, prod.Code())
		assert.True(t, prod.Price().Equal(prodPrice))
		assert.Nil(t, prod.Category())
		assert.NotNil(t, prod.Variants())
		assert.Empty(t, prod.Variants())
	})

	t.Run("product equality", func(t *testing.T) {
		prod1 := catalog.RestoreProduct(prodCode, prodPrice)
		prod2 := catalog.RestoreProduct(prodCode, prodPrice)
		assert.True(t, prod1.Equal(prod2))
	})

	t.Run("restore product with category", func(t *testing.T) {
		catCode, err := catalog.NewCategoryCode(1)
		require.NoError(t, err)
		cat, err := catalog.NewCategory(catCode, "Shoes")
		require.NoError(t, err)

		prod := catalog.RestoreProduct(prodCode, prodPrice, catalog.ProductWithCategory(cat))
		assert.NotNil(t, prod.Category())
		assert.True(t, prod.Category().Code().Equal(catCode))
	})

	t.Run("restore product with variants", func(t *testing.T) {

		variant1Price, err := price.Parse("5.99")
		require.NoError(t, err)
		variant1 := catalog.RestoreVariant("Variant 1", "VAR1", &variant1Price)
		variant2Price, err := price.Parse("7.99")
		require.NoError(t, err)
		variant2 := catalog.RestoreVariant("Variant 2", "VAR2", &variant2Price)

		prod := catalog.RestoreProduct(prodCode, prodPrice, catalog.ProductWithVariants(variant1, variant2))
		assert.Len(t, prod.Variants(), 2)

		variant := prod.Variants()[0]
		assert.Equal(t, "Variant 1", variant.Name())
		assert.Equal(t, "VAR1", variant.SKU())
		assert.True(t, variant1.Price().Equal(variant1Price))

		variant = prod.Variants()[1]
		assert.Equal(t, "Variant 2", variant.Name())
		assert.Equal(t, "VAR2", variant.SKU())
		assert.True(t, variant2.Price().Equal(variant2Price))
	})

	t.Run("fill variant with null price", func(t *testing.T) {

		variant1Price, err := price.Parse("5.99")
		require.NoError(t, err)
		variant1 := catalog.RestoreVariant("Variant 1", "VAR1", &variant1Price)
		variant2 := catalog.RestoreVariant("Variant 2", "VAR2", nil)

		prod := catalog.RestoreProduct(prodCode, prodPrice, catalog.ProductWithVariants(variant1, variant2))
		assert.Len(t, prod.Variants(), 2)

		assert.True(t, variant1.Price().Equal(variant1Price))
		assert.True(t, variant2.Price().Equal(prodPrice))
	})

}
