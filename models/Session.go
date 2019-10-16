package models

import (
	"math/rand"
	"time"
)

// Session структура сессии
type Session struct {
	SessID string
	UserID int
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
func CreateSession(userID int) (*Session, error) {
	sess := new(Session)
	sess.UserID = userID
	sess.SessID = randStringRunes(20)

	err := db.QueryRow("INSERT INTO sessions (sessid, user_id) values ($1, $2) RETURNING sessid", sess.SessID, userID).Scan(&sess.SessID)

	return sess, err
}

// CloseSession закрываем сессию
func CloseSession(sessID string) error {
	_, err := db.Exec("UPDATE sessions SET expires = CURRENT_TIMESTAMP WHERE sessid = $1", sessID)
	return err
}
