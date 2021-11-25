package error

import (
	"reflect"
)

// ErrCode type
type ErrCode string

const (
	// ErrCodeDefault ...
	ErrCodeDefault ErrCode = "default"
)

// Error for managed errors
type Error struct {
	Code    ErrCode     `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// Error interface
func (e *Error) Error() string {
	return e.Message
}

// Is interface
func (e *Error) Is(err error) bool {
	return reflect.TypeOf(e) == reflect.TypeOf(err)
}
