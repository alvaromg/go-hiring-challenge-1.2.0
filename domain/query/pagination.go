package query

type Pagination struct {
	page     int
	pageSize int
}

func NewPagination(page, pageSize int) Pagination {
	return Pagination{
		page:     page,
		pageSize: pageSize,
	}
}

func (p *Pagination) Page() int {
	return p.page
}

func (p *Pagination) PageSize() int {
	return p.pageSize
}
