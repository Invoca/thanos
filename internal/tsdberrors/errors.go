// Copyright (c) The Thanos Authors.
// Licensed under the Apache License 2.0.

// Package tsdberrors provides multi-error helpers removed from Prometheus 3.13
// (prometheus/prometheus#17768). Thanos still uses these until upstream adapts.
package tsdberrors

import (
	"bytes"
	"errors"
	"fmt"
	"io"
)

type multiError []error

// NewMulti returns multiError with provided errors added if not nil.
func NewMulti(errs ...error) multiError {
	m := multiError{}
	m.Add(errs...)
	return m
}

// Add adds single or many errors to the error list. Each error is added only if not nil.
func (es *multiError) Add(errs ...error) {
	for _, err := range errs {
		if err == nil {
			continue
		}
		var merr nonNilMultiError
		if errors.As(err, &merr) {
			*es = append(*es, merr.errs...)
			continue
		}
		*es = append(*es, err)
	}
}

// Err returns the error list as an error or nil if it is empty.
func (es multiError) Err() error {
	if len(es) == 0 {
		return nil
	}
	return nonNilMultiError{errs: es}
}

type nonNilMultiError struct {
	errs multiError
}

func (es nonNilMultiError) Error() string {
	var buf bytes.Buffer

	if len(es.errs) > 1 {
		fmt.Fprintf(&buf, "%d errors: ", len(es.errs))
	}

	for i, err := range es.errs {
		if i != 0 {
			buf.WriteString("; ")
		}
		buf.WriteString(err.Error())
	}

	return buf.String()
}

func (es nonNilMultiError) Is(target error) bool {
	for _, err := range es.errs {
		if errors.Is(err, target) {
			return true
		}
	}
	return false
}

func (es nonNilMultiError) Unwrap() []error {
	return es.errs
}

// CloseAll closes all given closers while recording error in MultiError.
func CloseAll(cs []io.Closer) error {
	errs := NewMulti()
	for _, c := range cs {
		errs.Add(c.Close())
	}
	return errs.Err()
}
