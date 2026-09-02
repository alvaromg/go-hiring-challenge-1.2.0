package catalog_test

import (
	"net/http"
	"testing"

	"github.com/mytheresa/go-hiring-challenge/test/helper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testCreateCategories(t *testing.T, router http.Handler) {
	t.Run("create categories", func(t *testing.T) {

		body := map[string]any{
			"categories": []map[string]string{
				{"code": "CAT901", "name": "test category 1"},
				{"code": "CAT902", "name": "test category 2"},
			},
		}

		rec := helper.DoRequest(t, router, http.MethodPost, "/categories", body)
		require.Equal(t, http.StatusCreated, rec.Code)

		// assert response

		var out helper.ResponseDTO[categoriesListDTO]
		helper.DecodeBody(t, rec, &out)

		assert.Len(t, out.Data.Categories, 2)

		cat := out.Data.Categories[0]
		assert.Equal(t, "CAT901", cat.Code)
		assert.Equal(t, "test category 1", cat.Name)

		cat = out.Data.Categories[1]
		assert.Equal(t, "CAT902", cat.Code)
		assert.Equal(t, "test category 2", cat.Name)

		// assert categories are persisted

		assertCategoriesExist(t, router, []string{"CAT901", "CAT902"})
	})

	t.Run("create category with invalid code", func(t *testing.T) {
		body := map[string]any{
			"categories": []map[string]string{
				{"code": "invalid", "name": "test category 1"},
			},
		}

		rec := helper.DoRequest(t, router, http.MethodPost, "/categories", body)
		require.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("create category with existent code", func(t *testing.T) {
		body := map[string]any{
			"categories": []map[string]string{
				{"code": "CAT901", "name": "test category 1"},
			},
		}

		rec := helper.DoRequest(t, router, http.MethodPost, "/categories", body)
		require.Equal(t, http.StatusConflict, rec.Code)
	})

	t.Run("create categories only applies as a unit", func(t *testing.T) {
		body := map[string]any{
			"categories": []map[string]string{
				{"code": "CAT901", "name": "existent category"},
				{"code": "CAT999", "name": "not existent category"},
			},
		}

		rec := helper.DoRequest(t, router, http.MethodPost, "/categories", body)
		require.Equal(t, http.StatusConflict, rec.Code)

		// assert no categories are persisted

		assertCategoriesDontExist(t, router, []string{"CAT999"})
	})
}

func assertCategoriesExist(t *testing.T, router http.Handler, codes []string) {
	rec := helper.DoRequest(t, router, http.MethodGet, "/categories", nil)
	require.Equal(t, http.StatusOK, rec.Code)

	var out helper.ListResponseDTO[categoriesListDTO]
	helper.DecodeBody(t, rec, &out)

	foundCats := map[string]struct{}{}
	for _, c := range out.Data.Categories {
		foundCats[c.Code] = struct{}{}
	}

	for _, c := range codes {
		assert.Contains(t, foundCats, c)
	}
}

func assertCategoriesDontExist(t *testing.T, router http.Handler, codes []string) {
	rec := helper.DoRequest(t, router, http.MethodGet, "/categories", nil)
	require.Equal(t, http.StatusOK, rec.Code)

	var out helper.ListResponseDTO[categoriesListDTO]
	helper.DecodeBody(t, rec, &out)

	foundCats := map[string]struct{}{}
	for _, c := range out.Data.Categories {
		foundCats[c.Code] = struct{}{}
	}

	for _, c := range codes {
		assert.NotContains(t, foundCats, c)
	}
}
