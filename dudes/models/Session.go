package models

import (
	"log"
	"math/rand"
	"time"
)

// Session структура сессии
type Session struct {
	SessID string
	UserID int64
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

// CreateSession создаем сессию и пишем в базу
func CreateSession(userID int64) (*Session, error) {
	sess := new(Session)
	sess.UserID = userID
	sess.SessID = randStringRunes(20)

	var curTime time.Time = time.Now()
	log.Println(curTime)
	_, err := dbM.Exec(
		"INSERT INTO sessions (sessid, user_id, expires) values ($1, $2, $3::timestamp + INTERVAL '24 hours')",
		sess.SessID, userID, curTime)

	return sess, err
}

// CloseSession закрываем сессию
func CloseSession(sessID string) error {

	var curTime time.Time = time.Now()

	_, err := dbM.Exec(
		"UPDATE sessions SET expires = $1 WHERE sessid = $2",
		curTime, sessID)
	return err
}
