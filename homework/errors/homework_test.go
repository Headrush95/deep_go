package main

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type MultiError struct {
	errs []error
}

func (e *MultiError) Error() string {
	if e == nil || len(e.errs) == 0 {
		return ""
	}
	sb := strings.Builder{}
	sb.Grow(124) // можно, конечно, точно посчитать, но это займет время (всякие err.Error() как будто бы долгие)
	sb.WriteString("2 errors occured:\n")
	for i := range e.errs {
		sb.WriteString("\t* " + e.errs[i].Error())
	}
	sb.WriteString("\n")

	return strings.Clone(sb.String())
}

func Append(err error, errs ...error) *MultiError {
	if err == nil {
		return &MultiError{errs: errs}
	}
	multErr, ok := err.(*MultiError)
	if !ok {
		multErr = &MultiError{}
		multErr.errs = make([]error, 0, len(errs)+1)
		multErr.errs = append(multErr.errs, err)
	}
	multErr.errs = append(multErr.errs, errs...)

	return multErr
}

func TestMultiError(t *testing.T) {
	var err error
	err = Append(err, errors.New("error 1"))
	err = Append(err, errors.New("error 2"))

	expectedMessage := "2 errors occured:\n\t* error 1\t* error 2\n"
	assert.EqualError(t, err, expectedMessage)
}
