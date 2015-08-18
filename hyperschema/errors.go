package hyperschema

import (
	"fmt"

	"github.com/davecgh/go-spew/spew"
)

const (
	ErrorCodeNotSupportedFile = 0 + iota
	ErrorCodeExtraSlash
	ErrorCodeDuplicated
	ErrorCodeNoID
	ErrorCodeProperty
	ErrorCodePropertyIncorrect
	ErrorCodeInvalidRefs
)

type ErrorPropertyType string

func (p ErrorPropertyType) String() string {
	return string(p)
}

const (
	ErrorPropertyTypeNil        = ErrorPropertyType("is nil")
	ErrorPropertyTypeEmpty      = ErrorPropertyType("is empty")
	ErrorPropertyTypeConflicted = ErrorPropertyType("are conflicted")
)

func ErrorNotSupportedFile(internal error) *Error {
	return &Error{Code: ErrorCodeNotSupportedFile, Message: "not supported file", Internal: internal}
}
func ErrorExtraSlash(id string) *Error {
	return &Error{Code: ErrorCodeExtraSlash, Message: fmt.Sprintf("extra slash in ID '%s'", id)}
}
func ErrorDuplicated(id string) *Error {
	return &Error{Code: ErrorCodeDuplicated, Message: fmt.Sprintf("duplicated ID '%s'", id)}
}
func ErrorNoID() *Error {
	return &Error{Code: ErrorCodeNoID, Message: "ID is empty"}
}
func ErrorProperty(id string, propertyName string, errorType ErrorPropertyType) *Error {
	return &Error{Code: ErrorCodeProperty, Message: fmt.Sprintf("%s in '%s' %s", propertyName, id, errorType.String())}
}
func ErrorPropertyIncorrect(id string, propertyName string, value interface{}) *Error {
	return &Error{Code: ErrorCodePropertyIncorrect, Message: fmt.Sprintf("%s in '%s', %s is incorrect", propertyName, id, spew.Sdump(value))}
}
func ErrorInvalidRefs(id string, ref string) *Error {
	return &Error{Code: ErrorCodeInvalidRefs, Message: fmt.Sprintf("referred from '%s' to %s but it's not found", id, ref)}
}

func addError(e **Error, new *Error) {
	if new == nil {
		return
	}
	if e == nil {
		panic("invalid call for hyperschema.AddError")
	}
	if *e == nil {
		*e = new
	} else {
		new.Prev = *e
		*e = new
	}
}

type Error struct {
	Code     int
	Message  string
	Internal error
	Prev     *Error
}

func (e *Error) Error() string {
	var message string
	if e.Prev != nil {
		message = e.Prev.Error() + "\n"
	}
	if e.Internal == nil {
		return message + fmt.Sprintf("%s(%d)", e.Message, e.Code)
	}
	return message + fmt.Sprintf("%s(%d);internal:%s", e.Message, e.Code, e.Internal.Error())
}
