package database

import (
	"database/sql"
	"testing"
)

func TestOpenConnection(t *testing.T) {
	db := openTestConnection(t)
	if err := db.Ping(); err != nil {
		t.Fatal("Unexpected error during ping database", err)
	}
	if err := db.Close(); err != nil {
		t.Fatal("Unexpected error during close database", err)
	}
}

func openTestConnection(t *testing.T) *sql.DB {
	db, err := OpenConnection("postgresql://user:password@localhost/funny-bot?sslmode=disable")
	if err != nil {
		t.Fatal("Unexpected error during open database connection", err)
	}
	return db
}
