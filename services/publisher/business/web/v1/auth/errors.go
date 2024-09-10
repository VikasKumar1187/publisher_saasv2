package auth

import (
	"errors"
	"fmt"
)



type AuthenticationError struct {
	msg string
}


type AuthorizationError struct {
	msg string
}


func NewAuthenticationError(format string, args ...any) error {
	return &AuthenticationError{
		msg: fmt.Sprintf(format, args...),
	}
}


func NewAuthorizationError(format string, args ...any) error {
	return &AuthorizationError{
		msg: fmt.Sprintf(format, args...),
	}
}


func (ae *AuthenticationError) Error() string {
	return ae.msg
}


func (ae *AuthorizationError) Error() string {
	return ae.msg
}


func IsAuthenticationError(err error) bool {
	var ae *AuthenticationError
	return errors.As(err, &ae)
}


func IsAuthorizationError(err error) bool {
	var ae *AuthorizationError
	return errors.As(err, &ae)
}

