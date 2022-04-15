package db

import (
	"database/sql"
	"os"

	_ "github.com/lib/pq"
)

func getConnection() (map[string]*sql.DB, error) {
	defaultDBUri := os.Getenv("DEFAULT_DB")
	defaultDB, _ := sql.Open("postgres", defaultDBUri)

	escortProfileDBUri := os.Getenv("ESCORT_PROFILE_DB")
	escortProfileDB, _ := sql.Open("postgres", escortProfileDBUri)

	return map[string]*sql.DB{
		"default":       defaultDB,
		"escortProfile": escortProfileDB,
	}, nil
}
