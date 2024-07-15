package db

import (
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
)

const (
	dbPath = "internal/db/sqlite/home.db"
)

var dao *sql.DB

func InitDB() error {
	var err error
    dao, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}

	return nil
}

func CloseDB() {
    if dao != nil {
        dao.Close()
    }
}
