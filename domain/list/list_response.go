package list

type ListResponse[T any] struct {
	items []T
	total uint
}

func NewListResponse[T any]() ListResponse[T] {
	return ListResponse[T]{
		items: []T{},
		total: 0,
	}
}

func (r ListResponse[T]) Items() []T {
	return r.items
}

func (r ListResponse[T]) Total() uint {
	return r.total
}

func (r *ListResponse[T]) AddItems(items ...T) {
	r.items = append(r.items, items...)
}

func (r *ListResponse[T]) SetTotal(total uint) {
	r.total = total
}
