package models

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/hl/utils"
)

// User model
type User struct {
	ID        int64
	Email     string
	Firstname string
	Lastname  string
}

// CreateUser функция создает запись в бд и возвращает структуру юзера с id новой записи
func CreateUser(firstname string, lastname string, email string, password string) (*User, error) {

	hashedPsswd := utils.HashAndSalt([]byte(password))

	// userID that will be returned after SQL insertion

	result, err := dbM.Exec("INSERT INTO users (firstname, lastname, email, password) values (?, ?, ?, ?)",
		firstname, lastname, email, hashedPsswd)

	usr := new(User)
	log.Println(0)
	usr.ID, err = result.LastInsertId()
	log.Println(1)
	log.Println(usr.ID)
	log.Println(err)

	usr.Email = email
	usr.Firstname = firstname
	usr.Lastname = lastname

	return usr, err

}

// UsersList вернёт срез структур User сформированных из строк таблицы  Users
func UsersList() ([]*User, error) {

	rows, err := dbS.Query("SELECT id, email, firstname, lastname FROM users")

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
	row := dbS.QueryRow("SELECT id, email, firstname, lastname FROM users WHERE id = ?", userID)
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

	row := dbS.QueryRow("SELECT id, firstname, lastname, email, password FROM users WHERE email = ?", email)
	err := row.Scan(&usr.ID, &usr.Firstname, &usr.Lastname, &usr.Email, &hashedPswrd)
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

// GetCurrentUser взять текущего юзера из сессии
func GetCurrentUser(sessID string) (*User, error) {
	usr := new(User)
	row := dbS.QueryRow(`SELECT id, firstname, lastname, email FROM users 
						 WHERE id = (SELECT s.user_id FROM sessions s WHERE 
									 s.sessid = ? 
									 AND s.expires >= CURRENT_TIMESTAMP
									 ORDER BY expires DESC LIMIT 1)`, sessID)
	err := row.Scan(&usr.ID, &usr.Firstname, &usr.Lastname, &usr.Email)
	if err != nil && err == sql.ErrNoRows {
		return nil, errors.New("no sessions")
	}
	return usr, err
}

// SearchUsers взять текущего юзера из сессии
func SearchUsers(query string) ([]*User, error) {

	users := make([]*User, 0)

	firstname, lastname := "", ""
	s := strings.Split(query, " ")

	log.Println(len(s))

	if len(s) < 1 && len(s) > 2 {
		return users, nil
	} else if len(s) == 1 {
		firstname = s[0] + "%"
		lastname = firstname
		log.Println("horosho")

		log.Println(firstname, lastname)
		query = `(SELECT id, email, firstname, lastname FROM users WHERE firstname LIKE ? ORDER BY id)
				  UNION ALL 
				 (SELECT id, email, firstname, lastname FROM users WHERE lastname LIKE ? ORDER BY id);`
	} else {
		log.Println("ploho")

		firstname, lastname = s[0]+"%", s[1]+"%"
		query = `SELECT id, email, firstname, lastname FROM users WHERE firstname LIKE ?
				 AND lastname LIKE ? ORDER BY id;`
	}

	rows, err := dbS.Query(query, firstname, lastname)

	if err != nil {
		log.Println("serach user")
		log.Println(err)
		// log.Fatal(err)
	}

	defer rows.Close()

	for rows.Next() {
		usr := new(User)
		err := rows.Scan(&usr.ID, &usr.Email, &usr.Firstname, &usr.Lastname)
		if err != nil {
			log.Println(err)
			// log.Fatal(err)
		}
		users = append(users, usr)
	}
	if err = rows.Err(); err != nil {
		log.Println(err)
		// log.Fatal(err)
	}
	return users, err
}
