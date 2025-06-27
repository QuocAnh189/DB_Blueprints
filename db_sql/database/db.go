package database

import (
	"database/sql"
	"db_blueprints/config"
	"fmt"
)

type IDatabase interface {
}

type Database struct {
	db *sql.DB
}

func NewDatabase(config *config.Config) (*Database, error) {
	// 1. Construct the Data Source Name (DSN) string from the config
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.DB_USER,
		config.DB_PASSWORD,
		config.DB_HOST,
		config.DB_PORT,
		config.DB_NAME,
	)

	// 2. Open the database connection
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		// 3. If connection fails, return the error
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// 4. On success, wrap the connection in your struct and return it
	dbSQL := &Database{
		db: db,
	}

	fmt.Println("Successfully connected to the database!")
	return dbSQL, nil
}
