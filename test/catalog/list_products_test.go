package catalog_test

import (
	"net/http"
	"testing"

	"github.com/mytheresa/go-hiring-challenge/test/helper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testListProducts(t *testing.T, router http.Handler) {
	t.Run("list products", func(t *testing.T) {
		rec := helper.DoRequest(t, router, http.MethodGet, "/v1/catalog", nil)
		require.Equal(t, http.StatusOK, rec.Code)

		var out helper.ListResponseDTO[productsListDTO]
		helper.DecodeBody(t, rec, &out)

		assert.NotEmpty(t, out.Metadata.OperationId)
		assert.EqualValues(t, 8, out.Metadata.TotalCount)
		assert.Equal(t, 1, out.Metadata.Page)
		assert.Equal(t, 10, out.Metadata.PageSize)
		assert.Len(t, out.Data.Products, 8)
	})

	t.Run("list products filtered by price less than", func(t *testing.T) {
		rec := helper.DoRequest(t, router, http.MethodGet, "/v1/catalog?filter_price_lt=10", nil)
		require.Equal(t, http.StatusOK, rec.Code)

		var out helper.ListResponseDTO[productsListDTO]
		helper.DecodeBody(t, rec, &out)

		var codes []string
		for _, p := range out.Data.Products {
			codes = append(codes, p.Code)
		}
		assert.ElementsMatch(t, []string{"PROD003", "PROD006", "PROD008"}, codes)
	})

	t.Run("list products filtered by category", func(t *testing.T) {
		rec := helper.DoRequest(t, router, http.MethodGet, "/v1/catalog?filter_category_eq=CAT001", nil)
		require.Equal(t, http.StatusOK, rec.Code)

		var out helper.ListResponseDTO[productsListDTO]
		helper.DecodeBody(t, rec, &out)

		assert.EqualValues(t, 3, out.Metadata.TotalCount)

		var codes []string
		for _, p := range out.Data.Products {
			codes = append(codes, p.Code)
			require.NotNil(t, p.Category)
			assert.Equal(t, "CAT001", p.Category.Code)
		}
		assert.ElementsMatch(t, []string{"PROD001", "PROD004", "PROD007"}, codes)
	})

	t.Run("list products pagination", func(t *testing.T) {
		var out helper.ListResponseDTO[productsListDTO]

		// page 1
		rec := helper.DoRequest(t, router, http.MethodGet, "/v1/catalog?pageSize=3&page=1", nil)
		require.Equal(t, http.StatusOK, rec.Code)

		helper.DecodeBody(t, rec, &out)

		assert.EqualValues(t, 8, out.Metadata.TotalCount)
		assert.Equal(t, 1, out.Metadata.Page)
		assert.Equal(t, 3, out.Metadata.PageSize)
		assert.Len(t, out.Data.Products, 3)

		// page 2
		rec = helper.DoRequest(t, router, http.MethodGet, "/v1/catalog?pageSize=3&page=2", nil)
		require.Equal(t, http.StatusOK, rec.Code)

		helper.DecodeBody(t, rec, &out)

		assert.EqualValues(t, 8, out.Metadata.TotalCount)
		assert.Equal(t, 2, out.Metadata.Page)
		assert.Equal(t, 3, out.Metadata.PageSize)
		assert.Len(t, out.Data.Products, 3)

		// page 3
		rec = helper.DoRequest(t, router, http.MethodGet, "/v1/catalog?pageSize=3&page=3", nil)
		require.Equal(t, http.StatusOK, rec.Code)

		helper.DecodeBody(t, rec, &out)

		assert.EqualValues(t, 8, out.Metadata.TotalCount)
		assert.Equal(t, 3, out.Metadata.Page)
		assert.Equal(t, 3, out.Metadata.PageSize)
		assert.Len(t, out.Data.Products, 2)

		// page 4
		rec = helper.DoRequest(t, router, http.MethodGet, "/v1/catalog?pageSize=3&page=4", nil)
		require.Equal(t, http.StatusOK, rec.Code)

		helper.DecodeBody(t, rec, &out)

		assert.EqualValues(t, 8, out.Metadata.TotalCount)
		assert.Equal(t, 4, out.Metadata.Page)
		assert.Equal(t, 3, out.Metadata.PageSize)
		assert.Len(t, out.Data.Products, 0)

	})

	t.Run("list products filtered by category that does not exist", func(t *testing.T) {
		rec := helper.DoRequest(t, router, http.MethodGet, "/v1/catalog?filter_category_eq=CAT999", nil)
		require.Equal(t, http.StatusNotFound, rec.Code)

		var out helper.ErrorResponseDTO
		helper.DecodeBody(t, rec, &out)

		assert.Equal(t, "not found error: category \"CAT999\" not found", out.Error.Message)
	})
}
