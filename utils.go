package rockgo

import (
	"errors"
	"fmt"
)

func NewError(v ...interface{}) error {
	return errors.New(fmt.Sprint(v...))
}
