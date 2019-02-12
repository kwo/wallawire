package model

import (
	"errors"
)

type NotFoundError struct {
	error
}

func (z *NotFoundError) NotFound() {}

func NewNotFoundError(msg string) *NotFoundError {
	return &NotFoundError{error: errors.New(msg)}
}

func IsNotFoundError(err error) bool {
	type NotFound interface {
		NotFound()
	}
	if _, ok := err.(NotFound); ok {
		return true
	}
	return false
}

type ValidationError struct {
	error
}

func (z *ValidationError) Invalid() {}

func NewValidationError(msg string) *ValidationError {
	return &ValidationError{error: errors.New(msg)}
}

func IsValidationError(err error) bool {
	type Invalid interface {
		Invalid()
	}
	if _, ok := err.(Invalid); ok {
		return true
	}
	return false
}
