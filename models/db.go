package models

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("postgres", "postgres://go_app:go_pass@localhost/go_app_db")
	if err != nil {
		log.Fatal(err)
	}
}
