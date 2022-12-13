package entity

import (
	"errors"
	"fmt"
)

// Ошибки для типа Metric
var (
	ErrParseValue   = errors.New("parsing error")
	ErrUnknownType  = errors.New("unknown type")
	ErrHashMismatch = errors.New("hash mismatch")
)

type (
	valueError struct {
		err   error
		name  string
		value string
	}

	typeError struct {
		err   error
		name  string
		tName string
	}

	hashError struct {
		err  error
		name string
	}
)

func (ve valueError) Error() string {
	return fmt.Sprintf("value %q error: %v - %v", ve.name, ve.value, ve.err)
}

func (ve valueError) Unwrap() error {
	return ve.err
}

func (te typeError) Error() string {
	return fmt.Sprintf("type %q error: %v - %v", te.name, te.tName, te.err)
}

func (te typeError) Unwrap() error {
	return te.err
}

func (he hashError) Error() string {
	return fmt.Sprintf("hash %q error: %v", he.name, he.err)
}

func (he hashError) Unwrap() error {
	return he.err
}
