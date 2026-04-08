package domain

import "errors"

var (
	ErrNotFound            = errors.New("resource not found")
	ErrConflict            = errors.New("resource already exists")
	ErrInvalidInput        = errors.New("invalid input")
	ErrTenantInactive      = errors.New("tenant is inactive")
	ErrApplicationInactive = errors.New("application is inactive")
	ErrForbidden           = errors.New("forbidden")
	ErrNotLinked           = errors.New("application is not linked to tenant")
)
