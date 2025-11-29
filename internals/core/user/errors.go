package user

import "errors"

var (
	ErrSomethingWentWrong      error = errors.New("something went wrong")
	ErrUserAlreadyRegistered   error = errors.New("user already registered")
	ErrInvalidUserDetail       error = errors.New("user detail is invalid")
	ErrUserNotFound            error = errors.New("user not found")
	ErrUnexpectedSigningMethod error = errors.New("unexpected token signin method")
	ErrInvalidRefreshToken     error = errors.New("invalid refresh token")
	ErrGettingDataFromDB       error = errors.New("error getting data from DB")
)

type ErrorStruct struct {
	Error    error
	ErrorMsg string
}
