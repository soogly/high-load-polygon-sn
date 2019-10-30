package models

import (
	"database/sql"
	"log"

	// _ "github.com/lib/pq"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("mysql", "go_app:go_pass@/go_app_db")
	if err != nil {
		log.Fatal(err)
	}
}
