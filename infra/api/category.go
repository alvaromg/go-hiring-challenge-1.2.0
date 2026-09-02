package api

import (
	"github.com/mytheresa/go-hiring-challenge/domain/catalog"
)

type category struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

func encodeCategoryResponse(c *catalog.Category) category {
	return category{
		Code: c.Code().String(),
		Name: c.Name(),
	}
}
