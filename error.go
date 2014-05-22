package jsonmap

import (
	"strings"
)

type Error []string

func (e *Error) Add(err error) {
	*e = append(*e, err.Error())
}

func (e Error) Empty() bool {
	return len(e) == 0
}

func (e Error) Error() string {
	return strings.Join(e, ", ")
}
