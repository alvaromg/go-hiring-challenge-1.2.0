package catalog

type Category struct {
	code string
	name string
}

func (c *Category) Code() string {
	return c.code
}

func (c *Category) Name() string {
	return c.name
}

func RestoreCategory(code, name string) *Category {
	category := &Category{
		code: code,
		name: name,
	}

	return category
}
