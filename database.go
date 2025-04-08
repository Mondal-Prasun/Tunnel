package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type SqlDb struct {
	Db *sql.DB
}

func initDataBase() (*SqlDb, error) {

	if _, err := os.Stat(DATABASE_DIRECTORY); err != nil {
		err = os.Mkdir(DATABASE_DIRECTORY, os.ModeDir)

		if err != nil {
			return nil, err
		}

	}

	db, err := sql.Open(DATABASE_DRIVER_NAME, DATABASE_PATH)

	if err != nil {
		log.Println("Cannot open sql driver:", err.Error())
		return nil, err
	}

	userData := `CREATE TABLE user(
    id TEXT NOT NULL,
    name TEXT NOT NULL,
    password TEXT NOT NULL,
    image TEXT NOT NULL);`

	if _, err := db.Exec(userData); err != nil {
		return nil, err
	}

	return &SqlDb{
		Db: db,
	}, nil

}
