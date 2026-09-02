package query

// Filter represents a single filter condition
type Filter struct {
	field    string
	operator Operator
	value    any
}

func NewFilter(field string, operator Operator, value any) Filter {
	return Filter{
		field:    field,
		operator: operator,
		value:    value,
	}
}
func (f *Filter) Field() string {
	return f.field
}
func (f *Filter) Operator() Operator {
	return f.operator
}
func (f *Filter) Value() any {
	return f.value
}
