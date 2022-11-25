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
	valueErr struct {
		name  string
		value string
		err   error
	}

	typeErr struct {
		name  string
		tName string
		err   error
	}

	hashErr struct {
		name string
		err  error
	}
)

func (ve valueErr) Error() string {
	return fmt.Sprintf("value %q error: %v - %v", ve.name, ve.value, ve.err)
}

func (ve valueErr) Unwrap() error {
	return ve.err
}

func (te typeErr) Error() string {
	return fmt.Sprintf("type %q error: %v - %v", te.name, te.tName, te.err)
}

func (te typeErr) Unwrap() error {
	return te.err
}

func (he hashErr) Error() string {
	return fmt.Sprintf("hash %q error: %v", he.name, he.err)
}

func (he hashErr) Unwrap() error {
	return he.err
}
