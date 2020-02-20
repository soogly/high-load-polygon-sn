package models

import (
	"database/sql"
	"log"
	"sort"
	"time"
)

// Dialog структура диалога
type Dialog struct {
	ID        int
	Firstname int
	Lastname  int
	LastMess  string
}

// Message структура диалога
type Message struct {
	ID        int
	Author    int
	Recipient int
	PubDate   time.Time
}

// GetDialogByUsers берем диалог по id юзеров
func GetDialogByUsers(members []int) *Dialog {

	var dialog = new(Dialog)
	sort.Ints(members)

	u1, u2 := members[0], members[1]

	row := dbS.QueryRow("SELECT id FROM dialog WHERE u1 = ? AND u2 = ?", u1, u2)

	err := row.Scan(&dialog.ID)

	if err == sql.ErrNoRows {
		return CreateDialog(u1, u2, "")
	}
	if err != nil {
		log.Println("Getting Dialog error")
		log.Fatal(err)
	}
	return dialog
}

// GetDialogs берем все диалоги пользователя
func GetDialogs(userID int) []*Dialog {

	rows, err := dbS.Query("SELECT d.id, d.last_mess, u.firstname, u.lastname FROM users u, dialogs d WHERE (u.id = d.u1 AND d.u2 = ? OR (u.id = d.u2 AND d.u1 = ?) AND u.id <> ?", userID)

	if err != nil {
		log.Println("Getting Dialogs List Error: ")
		log.Fatal(err)
	}

	defer rows.Close()

	var dialogs = make([]*Dialog, 0)
	for rows.Next() {
		dialog := new(Dialog)
		err := rows.Scan(&dialog.ID, &dialog.LastMess, &dialog.Firstname, &dialog.Lastname)
		if err != nil {
			log.Println("Scaning Dialogs List Error: ")
			log.Fatal(err)
		}
		dialogs = append(dialogs, dialog)
	}
	if err = rows.Err(); err != nil {
		log.Println("Some DB error when Getting Dialog List: ")
		log.Fatal(err)
	}
	return dialogs
}

// CreateDialog создаем диалог и пишем в базу
func CreateDialog(u1 int, u2 int, lastMess string) *Dialog {
	dialog := new(Dialog)
	dialog.LastMess = lastMess

	_, err := dbM.Exec(`INSERT INTO dialog (u1, u2, last_mess)
						values (?, ?, ?)`, u1, u2, lastMess)

	if err != nil {
		log.Println("Dialog Creation error: ")
		log.Fatal(err)
	}
	return dialog
}

// GetDialogMessages берем все сообщения диалога
func GetDialogMessages(dialogID int) []*Message {

	rows, err := dbS.Query("SELECT  id, author, recipient, pub_date FROM messages WHERE dialog = ?", dialogID)

	if err != nil {
		log.Println("Getting Dialogs Messages Error: ")
		log.Fatal(err)
	}

	defer rows.Close()

	var messages = make([]*Message, 0)
	for rows.Next() {
		mess := new(Message)
		err := rows.Scan(&mess.ID, &mess.Author, &mess.Recipient, &mess.PubDate)
		if err != nil {
			log.Println("Scaning Dialogs Messages List Error: ")
			log.Fatal(err)
		}
		messages = append(messages, mess)
	}
	if err = rows.Err(); err != nil {
		log.Println("Some DB error when Getting Dialogs Messages List: ")
		log.Fatal(err)
	}
	return messages
}
