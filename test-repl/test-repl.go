package main

import (
	"log"
	"time"

	"github.com/hl/models"

	// _ "github.com/lib/pq"
	_ "github.com/go-sql-driver/mysql"
)

// GenerateSessionss создаем сессию и пишем в базу
func main() {
	sessid := 677777
	log.Println(sessid)
	for sessid > 1 {
		sessid++
		log.Println(sessid)
		var curTime time.Time = time.Now()
		log.Println(curTime)
		res, err := models.DbM.Exec(`INSERT INTO sessions (sessid, user_id, expires)
									 values (?, ?, DATE_ADD(?, INTERVAL 1 DAY))`, sessid, 22, curTime)
		if err != nil {
			log.Fatal(err)
		}

		lid, err := res.LastInsertId()

		if err != nil {
			log.Fatal(err)
		}
		log.Println(lid)
	}
}
