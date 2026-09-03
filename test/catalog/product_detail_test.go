package catalog_test

import (
	"net/http"
	"testing"

	"github.com/mytheresa/go-hiring-challenge/test/helper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testProductDetail(t *testing.T, router http.Handler) {
	t.Run("product detail", func(t *testing.T) {
		rec := helper.DoRequest(t, router, http.MethodGet, "/v1/catalog/PROD001", nil)
		require.Equal(t, http.StatusOK, rec.Code)

		var out helper.ResponseDTO[productDTO]
		helper.DecodeBody(t, rec, &out)

		assert.Equal(t, "PROD001", out.Data.Code)
		assert.Equal(t, "10.99", out.Data.Price)
		assert.NotNil(t, out.Data.Category)
		assert.Equal(t, "CAT001", out.Data.Category.Code)
		assert.Equal(t, "Clothing", out.Data.Category.Name)

		assert.Len(t, out.Data.Variants, 3)

		variant := out.Data.Variants[0]
		assert.Equal(t, "Variant A", variant.Name)
		assert.Equal(t, "SKU001A", variant.SKU)
		assert.Equal(t, "11.99", *variant.Price)

		variant = out.Data.Variants[1]
		assert.Equal(t, "Variant B", variant.Name)
		assert.Equal(t, "SKU001B", variant.SKU)
		assert.Equal(t, "10.99", *variant.Price)

		variant = out.Data.Variants[2]
		assert.Equal(t, "Variant C", variant.Name)
		assert.Equal(t, "SKU001C", variant.SKU)
		assert.Equal(t, "10.99", *variant.Price)
	})

	t.Run("product detail not found", func(t *testing.T) {
		rec := helper.DoRequest(t, router, http.MethodGet, "/v1/catalog/PROD999", nil)
		require.Equal(t, http.StatusNotFound, rec.Code)

		var out helper.ErrorResponseDTO
		helper.DecodeBody(t, rec, &out)

		assert.Equal(t, "not found error: product \"PROD999\" not found", out.Error.Message)
	})
}
