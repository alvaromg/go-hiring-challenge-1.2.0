package catalog_test

import (
	"net/http"
	"testing"

	"github.com/mytheresa/go-hiring-challenge/test/helper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testListCategories(t *testing.T, router http.Handler) {
	t.Run("list categories", func(t *testing.T) {
		rec := helper.DoRequest(t, router, http.MethodGet, "/v1/categories", nil)
		require.Equal(t, http.StatusOK, rec.Code)

		var out helper.ListResponseDTO[categoriesListDTO]
		helper.DecodeBody(t, rec, &out)

		assert.EqualValues(t, 3, out.Metadata.TotalCount)
		assert.Equal(t, 1, out.Metadata.Page)
		assert.Equal(t, 10, out.Metadata.PageSize)
		assert.Len(t, out.Data.Categories, 3)
	})

	t.Run("list products pagination", func(t *testing.T) {
		var out helper.ListResponseDTO[categoriesListDTO]

		// page 1
		rec := helper.DoRequest(t, router, http.MethodGet, "/v1/categories?pageSize=2&page=1", nil)
		require.Equal(t, http.StatusOK, rec.Code)

		helper.DecodeBody(t, rec, &out)

		assert.EqualValues(t, 3, out.Metadata.TotalCount)
		assert.Equal(t, 1, out.Metadata.Page)
		assert.Equal(t, 2, out.Metadata.PageSize)
		assert.Len(t, out.Data.Categories, 2)

		// page 2
		rec = helper.DoRequest(t, router, http.MethodGet, "/v1/categories?pageSize=2&page=2", nil)
		require.Equal(t, http.StatusOK, rec.Code)

		helper.DecodeBody(t, rec, &out)

		assert.EqualValues(t, 3, out.Metadata.TotalCount)
		assert.Equal(t, 2, out.Metadata.Page)
		assert.Equal(t, 2, out.Metadata.PageSize)
		assert.Len(t, out.Data.Categories, 1)

		// page 3
		rec = helper.DoRequest(t, router, http.MethodGet, "/v1/categories?pageSize=2&page=3", nil)
		require.Equal(t, http.StatusOK, rec.Code)

		helper.DecodeBody(t, rec, &out)

		assert.EqualValues(t, 3, out.Metadata.TotalCount)
		assert.Equal(t, 3, out.Metadata.Page)
		assert.Equal(t, 2, out.Metadata.PageSize)
		assert.Len(t, out.Data.Categories, 0)

	})
}
