package query

type Sort struct {
	field string
	desc  bool
}

func NewSort(field string, desc bool) Sort {
	return Sort{
		field: field,
		desc:  desc,
	}
}

func (s *Sort) Field() string {
	if s == nil {
		return ""
	}
	return s.field
}

func (s *Sort) Desc() bool {
	if s == nil {
		return false
	}
	return s.desc
}
