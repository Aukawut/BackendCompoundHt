package handlers

import (
	"database/sql"

	"github.com/aukawut/BackendCompoundHt/config"
)

func OpenConnectDatabase() (*sql.DB, error) {
	connString := config.LoadDatabaseConfig()

	db, err := sql.Open("sqlserver", connString)

	return db, err
}
