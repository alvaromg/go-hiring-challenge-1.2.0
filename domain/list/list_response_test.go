package list_test

import (
	"testing"

	"github.com/mytheresa/go-hiring-challenge/domain/list"
	"github.com/stretchr/testify/assert"
)

func TestListResponse(t *testing.T) {
	t.Run("new list response starts empty", func(t *testing.T) {
		r := list.NewListResponse[string]()
		assert.Equal(t, []string{}, r.Items())
		assert.Equal(t, uint(0), r.Total())
	})

	t.Run("add items appends to the list", func(t *testing.T) {
		r := list.NewListResponse[int]()

		r.AddItems(1, 2, 3)
		assert.Equal(t, []int{1, 2, 3}, r.Items())

		r.AddItems(4)
		assert.Equal(t, []int{1, 2, 3, 4}, r.Items())
	})

	t.Run("set total overwrites the total", func(t *testing.T) {
		r := list.NewListResponse[int]()

		r.SetTotal(42)
		assert.Equal(t, uint(42), r.Total())

		r.SetTotal(0)
		assert.Equal(t, uint(0), r.Total())
	})
}
