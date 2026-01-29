package response

import "errors"

var (
	ErrInvalidJSON      = errors.New("invalid JSON format")
	ErrEmployeeExists   = errors.New("employee already exists")
	ErrEmployeeMissing  = errors.New("employee not found")
	ErrStorage          = errors.New("error while accessing data resource")
	ErrMethodNotAllowed = errors.New("method not allowed")
	ErrRouteNotFound    = errors.New("route not found")
	ErrMissingField     = errors.New("name, age and address are required")
)
