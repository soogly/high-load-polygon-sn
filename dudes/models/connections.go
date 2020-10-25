package models

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var dbM *sql.DB // Master db
var dbS *sql.DB // Slave db

func init() {
	var err error
	dbM, err = sql.Open("postgres", os.Getenv("DB_MASTER_URL"))
	if err != nil {
		log.Fatal(err)
	}

	dbS, err = sql.Open("postgres", os.Getenv("DB_SLAVE_URL"))
	if err != nil {
		log.Fatal(err)
	}
}
