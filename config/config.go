package config

import (
	"fmt"
	"os"
)

func LoadDatabaseConfig() string {

	// Set up the connection string for SQL Server
	connString := fmt.Sprintf("sqlserver://%s:%s@%s:1433?database=%s&encrypt=disable&connection+timeout=30",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_SERVER"),
		os.Getenv("DB_NAME"))

	return connString
}
