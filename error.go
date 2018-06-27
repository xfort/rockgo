package rockgo

import (
	"fmt"
)

type ErrorObj struct {
	Code int
	Tag  string
	Desc string
}

func NewError(tag string, v ...interface{}) *ErrorObj {
	return &ErrorObj{Tag: tag, Desc: fmt.Sprint(v...)}
}

func (errorobj *ErrorObj) Error() string {
	return fmt.Sprintf("%d_%s_%s", errorobj.Code, errorobj.Tag, errorobj.Desc)
}
