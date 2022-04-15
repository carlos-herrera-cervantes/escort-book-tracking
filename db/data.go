package db

import (
	"database/sql"
	"log"
)

var (
	data *Data
)

type Data struct {
	DB *sql.DB
}

func InitDB(db string) *Data {
	dbs, err := getConnection()

	if err != nil {
		log.Panic(err)
	}

	data = &Data{
		DB: dbs[db],
	}

	return data
}

func Close() error {
	if data == nil {
		return nil
	}

	return data.DB.Close()
}
