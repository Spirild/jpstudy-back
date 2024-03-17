package common

import (
	"errors"
)

type CommonError struct {
	ErrNoDatabaseService      error
	ErrInvalidDatabaseService error
	ErrNoThirdService         error
	ErrInvalidThirdService    error
}

var ErrorInstance = CommonError{
	ErrNoDatabaseService:      errors.New("could not find database service"),
	ErrInvalidDatabaseService: errors.New("database service doesn't work"),
	ErrNoThirdService:         errors.New("could not find third service"),
	ErrInvalidThirdService:    errors.New("third service doesn't work"),
}
