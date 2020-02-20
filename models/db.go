package models

import (
	"database/sql"
	"log"

	// _ "github.com/lib/pq"
	_ "github.com/go-sql-driver/mysql"
)

var dbM *sql.DB // Master db

// DbM Master
var DbM *sql.DB

var dbS *sql.DB // Slave db

func init() {
	var err error
	dbM, err = sql.Open("mysql", "go_app:go_pass@tcp(mysqlmaster:3306)/go_app_db")
	if err != nil {
		log.Fatal(err)
	}
	dbS, err = sql.Open("mysql", "go_app:go_pass@tcp(mysqlslave_1:3308)/go_app_db")
	if err != nil {
		log.Fatal(err)
	}

	DbM, err = sql.Open("mysql", "go_app:go_pass@tcp(172.23.0.4:3306)/go_app_db")
	if err != nil {
		log.Fatal(err)
	}
}
