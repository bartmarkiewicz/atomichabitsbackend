package util

import (
	"database/sql/driver"
	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
)

func NewMockDatabase() (*gorm.DB, sqlmock.Sqlmock, error) {
	database, mock, err := sqlmock.New()

	if err != nil {
		return nil, nil, err
	}

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: database}), &gorm.Config{})

	if err != nil {
		return nil, nil, err
	}

	return gormDB, mock, nil
}

type AnyTime struct{}

func (a AnyTime) Match(value driver.Value) bool {
	_, ok := value.(time.Time)
	return ok
}
