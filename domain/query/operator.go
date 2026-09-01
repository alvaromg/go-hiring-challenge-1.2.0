package query

import "slices"

// Operator represents comparison operators for criteria
type Operator string

const (
	Eq    Operator = "eq"
	Ne    Operator = "ne"
	Gt    Operator = "gt"
	Gte   Operator = "gte"
	Lt    Operator = "lt"
	Lte   Operator = "lte"
	In    Operator = "in"
	Nin   Operator = "nin"
	Like  Operator = "like"
	Is    Operator = "is"
	IsNot Operator = "is_not"
)

func ParseOperator(op string) (Operator, error) {
	o := Operator(op)
	if !o.IsValid() {
		return "", newErrInvalidQuery("invalid operator %s", op)
	}
	return o, nil
}

func (o Operator) IsValid() bool {
	return slices.Contains(AllOperators(), o)
}

func AllOperators() []Operator {
	return []Operator{Eq, Ne, Gt, Gte, Lt, Lte, In, Nin, Like, Is, IsNot}
}
