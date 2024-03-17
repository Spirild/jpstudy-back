package jpliteservice

import (
	"errors"
)

var (
	errNoDatabaseService      = errors.New("could not find database service")
	errInvalidDatabaseService = errors.New("database service doesn't work")
	errNoThirdService         = errors.New("could not find third service")
	errInvalidThirdService    = errors.New("third service doesn't work")
)
