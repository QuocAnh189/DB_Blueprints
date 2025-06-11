package gorm

import "gorm.io/gorm"

type IDatabase interface {
}

type Database struct {
	db *gorm.DB
}

func NewDatabase(uri string) (*Database, error) {
	return nil, nil
}
