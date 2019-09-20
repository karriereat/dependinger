package database

import (
	"database/sql"

	"log"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB
var err error

func Connect() *sql.DB {
	if db == nil {
		log.Println("No Connection creating one")
		db, err = sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/gohard")
		//TODO Change to config
		if err != nil {
			panic(err)
		}
	} else {
		log.Println("Returning existing connection")
	}

	return db
}

func Close() {
	db.Close()
}
