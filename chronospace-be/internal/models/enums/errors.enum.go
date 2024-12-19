package enums

import "errors"

type ErrorCode int

var (
	ErrInvalidContex         = errors.New("invalid context")
	ErrInvalidCredentials    = errors.New("invalid credentials")
	ErrPassword8Symbols      = errors.New("password must be at least 8 characters")
	ErrInvalidEmailFormat    = errors.New("invalid email format")
	ErrEmailAlreadyExists    = errors.New("email already exists")
	ErrUsernameAlreadyExists = errors.New("username already exists")
	ErrGeneratingToken       = errors.New("error generationg token")
	ErrStoringToken          = "error storing token: %s"
	ErrCleaningToken         = errors.New("error cleaning token")
	ErrDeletingUser          = errors.New("cannot to delete user")
	ErrListingUsers          = errors.New("listing users currently not possible")
	ErrUserNotFound          = errors.New("no user")

	ErrBookingInvalidInput     = errors.New("invalid input")
	ErrBookingInvalidDateRange = errors.New("invalid date range")
)
