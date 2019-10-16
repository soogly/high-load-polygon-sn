package models

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/hl/utils"
)

// User model
type User struct {
	ID        int
	Email     string
	Firstname string
	Lastname  string
}

// CreateUser функция создает запись в бн=д и возвращает структуру юзера с id новойзаписи
func CreateUser(firstname string, lastname string, email string, password string) (*User, error) {

	hashedPsswd := utils.HashAndSalt([]byte(password))
	fmt.Println(hashedPsswd)

	// userID that will be returned after SQL insertion
	var userID int

	err := db.QueryRow("INSERT INTO users (firstname, lastname, email, password) values ($1, $2, $3, $4 ) RETURNING id",
		firstname, lastname, email, hashedPsswd).Scan(&userID)

	usr := new(User)
	usr.ID = userID
	usr.Email = email
	usr.Firstname = firstname
	usr.Lastname = lastname

	fmt.Println(lastname)
	fmt.Println(err)
	fmt.Println(lastname)

	return usr, err

}

// UsersList вернёт срез структур User сформированных из строк таблицы  Users
func UsersList() ([]*User, error) {

	rows, err := db.Query("SELECT id, email, firstname, lastname FROM users")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	users := make([]*User, 0)
	for rows.Next() {
		usr := new(User)
		err := rows.Scan(&usr.ID, &usr.Email, &usr.Firstname, &usr.Lastname)
		if err != nil {
			log.Fatal(err)
		}
		users = append(users, usr)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
	return users, err
}

// UserProfile return User struct builded from User DB.table
func UserProfile(userID string) (*User, error) {

	usr := new(User)
	row := db.QueryRow("SELECT id, email, firstname, lastname FROM users WHERE id = $1", userID)
	err := row.Scan(&usr.ID, &usr.Email, &usr.Firstname, &usr.Lastname)

	if err != nil {
		log.Fatal(err)
	}

	return usr, err
}

// LoginUser check password and return User struct builded from User DB.table
func LoginUser(email string, password string) (*Session, error) {

	usr := new(User)
	var hashedPswrd string

	row := db.QueryRow("SELECT id, firstname, lastname, password FROM users WHERE email = $1", email)
	err := row.Scan(&usr.ID, &usr.Firstname, &usr.Lastname, &hashedPswrd)
	if err != nil && err == sql.ErrNoRows {
		return nil, errors.New("wrong email")
	}

	pswrdIsOk := utils.ComparePasswords(hashedPswrd, []byte(password))
	fmt.Println(pswrdIsOk)
	if pswrdIsOk == true {
		return CreateSession(usr.ID)
	}
	return nil, errors.New("wrong pass")
}
