package errors

import "errors"

var (
	ErrorDomainValidation   = errors.New("domain validation error")
	ErrorDuplicatedResource = errors.New("duplicated resource error")
	ErrorNotFound           = errors.New("not found error")
)
