// ./internal/api/apperror.go
package api

import(
	"runtime/debug"	
)

// appError type. will contain these field.
//
// Code: for the frontend logic. In case we require some consistent code. 
//
// Message: explain the error.
//
// Fileds: where the error is. Useful for validatoin error that might exist 
//
// Status: tfor backend to know. Not for frontend
//
// Err: This is the internal Error. Mostly for debugging and not exposing Backend Error to the frontend (the frontend should recieve expected error mostly)
//
// stack: to store the stack trace for error during development 
type AppError struct {	
	Code    string       `json:"code"`
	Message string       `json:"message"`
	Fields  []FieldError `json:"fields,omitempty"`

	Status int   `json:"-"`
	Err    error `json:"-"`
	Stack  []byte `json:"-"`
}

// this is for calling the error function on AppError (as itself is a error type. Which should return the message. As well as any internal error should also be exposed)
func (e *AppError) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}

	return e.Message
}

// to unwrap the error message inside the app Error
func (e *AppError) Unwrap() error {
	return e.Err
}

// to wrap more than one error type. This is useful to know where the exact error is when used fairly well.
func (e *AppError) Wrap(err error) *AppError {
	e.Err = err
	// capture stack at the real failure point
	e.Stack = debug.Stack()
	return e
}

// field Error for adding fields in case of validation and stuff
type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// for creating Errors  
func NewError(code, message string, status int) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Status:  status,
	}
}

func (e *AppError) AddField(field, message string) *AppError {
	e.Fields = append(e.Fields, FieldError{
		Field:   field,
		Message: message,
	})

	return e
}