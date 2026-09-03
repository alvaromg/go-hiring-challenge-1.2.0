package api

import (
	"net/http"
	"testing"

	"github.com/mytheresa/go-hiring-challenge/test/helper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testDefaultNotFound(t *testing.T, router http.Handler) {
	t.Run("default not found route", func(t *testing.T) {
		rec := helper.DoRequest(t, router, http.MethodGet, "/non-existent-resource?param1=value1", nil)
		require.Equal(t, http.StatusNotFound, rec.Code)

		var out helper.ErrorResponseDTO
		helper.DecodeBody(t, rec, &out)

		assert.Equal(t, "route \"/non-existent-resource\" not found", out.Error.Message)

	})

}
