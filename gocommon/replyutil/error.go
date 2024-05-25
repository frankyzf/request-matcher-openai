package replyutil

import (
	"errors"
	"fmt"
	"reflect"
)

// 12
type StaticFileAuthFailure struct {
	Message string
}

func (e StaticFileAuthFailure) Error() string {
	return fmt.Sprintf("static file auth failure %v", e.Message)
}

// 400
type AuthError struct {
	Message string
}

func (e AuthError) Error() string {
	return fmt.Sprintf("Auth Error:%v", e.Message)
}

// 401
type AuthExpireError struct {
	Message string
}

// 402
func (e AuthExpireError) Error() string {
	return fmt.Sprintf("Auth Expired:%v", e.Message)
}

type AuthPriviledgeError struct {
	Message string
}

func (e AuthPriviledgeError) Error() string {
	return fmt.Sprintf("Auth Privilege:%v", e.Message)
}

// 403
type AuthReloginError struct {
	Message string
}

func (e AuthReloginError) Error() string {
	return fmt.Sprintf("Auth Relogin:%v", e.Message)
}

func GetErrorCode(err error) int { //
	if reflect.TypeOf(err) == reflect.TypeOf(StaticFileAuthFailure{}) {
		return 12
	} else if reflect.TypeOf(err) == reflect.TypeOf(AuthError{}) {
		return 400
	} else if reflect.TypeOf(err) == reflect.TypeOf(AuthExpireError{}) {
		return 401
	} else if reflect.TypeOf(err) == reflect.TypeOf(AuthPriviledgeError{}) {
		return 402
	} else if reflect.TypeOf(err) == reflect.TypeOf(AuthReloginError{}) {
		return 403
	} else {
		return 500
	}
}

func GetError(code int, message string) error { //
	if code == 12 {
		return StaticFileAuthFailure{Message: message}
	} else if code == 400 {
		return AuthError{Message: message}
	} else if code == 401 {
		return AuthExpireError{Message: message}
	} else if code == 402 {
		return AuthPriviledgeError{Message: message}
	} else if code == 403 {
		return AuthReloginError{Message: message}
	} else {
		return errors.New(message)
	}
}
