package services

import "errors"

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrPasswordHashing    = errors.New("password hashing error")
	ErrCannotCreateUser   = errors.New("cannot create user")
	ErrCannotFindUser     = errors.New("cannot find user")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrCannotSignToken    = errors.New("cannot sign token")
	ErrCannotLoginUser    = errors.New("cannot login user")

	ErrCannotCreateAdvertisement = errors.New("cannot create advertisement")
	ErrCannotGetAdvertisements   = errors.New("cannot get advertisements")
)
