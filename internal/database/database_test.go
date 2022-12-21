package database

import "testing"

func TestOpenConnection(t *testing.T) {
	db, err := OpenConnection("postgresql://user:password@localhost/funny-bot?sslmode=disable")
	if err != nil {
		t.Fatal("Unexpected error during open database connection", err)
	}
	if err := db.Ping(); err != nil {
		t.Fatal("Unexpected error during ping database", err)
	}
	if err := db.Close(); err != nil {
		t.Fatal("Unexpected error during close database", err)
	}
}
