package models

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Thing struct {
	ID           int64      `json:"id"`
	Type         string     `json:"type"`
	Title        string     `json:"title"`
	CreationDate *time.Time `json:"creation_date"`
	Dude         int64      `json:"dude"`
	Comment      NullString `json:"comment"`
	Priority     int64      `json:"priority"`
	Duration     int32      `json:"duration"`
	WhenIt       *time.Time `json:"when_it"`
	Step         byte       `json:"step"`
	StartTime    NullString `json:"start_time"`
	StartsFrom   *time.Time `json:"starts_from"`
	OnlyIn       NullString `json:"only_in"`
	Done         NullBool   `json:"done"`
	BigDeal      NullInt64  `json:"big_deal"`
}

func GetUsersThingsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), 405)
		return
	}

	userID := r.FormValue("id")
	log.Println(userID, "userID")

	if userID == "" {
		http.Error(w, http.StatusText(404), 404)
		return
	}

	things := GetThings(userID)

	b, err := json.Marshal(things)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func GetThings(usrId string) []*Thing {
	// var conn pgxpool.Conn

	conn, err := dbpoolS.Acquire(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Release()

	rows, orr := conn.Query(context.Background(),
		`SELECT id, ttype, title, creation_date, dude, comment, priority, duration,
		 when_it, step, start_time, starts_from, only_in, done, big_deal
		 FROM things WHERE dude=$1`, usrId)

	if orr != nil {
		log.Println("GetThings3")
		log.Fatal(orr)
	}

	defer rows.Close()

	var things []*Thing
	for rows.Next() {
		thing := new(Thing)
		err := rows.Scan(&thing.ID, &thing.Type, &thing.Title, &thing.CreationDate,
			&thing.Dude, &thing.Comment, &thing.Priority, &thing.Duration, &thing.WhenIt,
			&thing.Step, &thing.StartTime, &thing.StartsFrom, &thing.OnlyIn,
			&thing.Done, &thing.BigDeal)
		if err != nil {
			log.Fatal(err)
		}
		things = append(things, thing)

	}
	if err = rows.Err(); err != nil {
		log.Println("GetThings4")
		log.Fatal(err)
	}
	return things
}

func CreateThingHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), 405)
		return
	}

	decoder := json.NewDecoder(r.Body)

	var thing Thing

	err := decoder.Decode(&thing)

	if err != nil { // bad request
		log.Println(r.Body)
		log.Println()
		log.Println(thing)
		log.Println()
		log.Println(thing.WhenIt)
		log.Println()
		log.Println(err)
		w.WriteHeader(400)
		return
	}

	userID := r.FormValue("userID")
	log.Println(userID, "userID")

	if userID == "" {
		http.Error(w, http.StatusText(404), 404)
		return
	}

	conn, err := dbpoolM.Acquire(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Release()

	row := conn.QueryRow(context.Background(),
		`INSERT INTO things (ttype, title, dude, comment, priority, duration,
		 when_it, step, start_time, starts_from, only_in, done, big_deal) VALUES 
		 ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) RETURNING id`,
		&thing.Type, &thing.Title, &userID, &thing.Comment, &thing.Priority,
		&thing.Duration, &thing.WhenIt,
		&thing.Step, &thing.StartTime, &thing.StartsFrom, &thing.OnlyIn,
		&thing.Done, &thing.BigDeal)

	var thing_id uint64
	err = row.Scan(&thing_id)

	if err != nil {
		log.Printf("Unable to INSERT: %v\n", err)
		w.WriteHeader(500)
		return
	}

	resp := make(map[string]string, 1)
	resp["thing_id"] = strconv.FormatUint(thing_id, 10)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Printf("Unable to encode json: %v\n", err)
		w.WriteHeader(500)
		return
	}
}
