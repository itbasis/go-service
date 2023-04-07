package db

import "errors"

var (
	ErrCreateInstance = errors.New("DB instance creation error")
	ErrDbMigration    = errors.New("database migration error")
	ErrConnectDB      = errors.New("database connection error")
)
