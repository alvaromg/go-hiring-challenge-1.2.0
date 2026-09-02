package api

type categoriesList struct {
	Categories []category `json:"categories"`
}

func wrapCategoriesResponse(categories []category) categoriesList {
	return categoriesList{Categories: categories}
}
