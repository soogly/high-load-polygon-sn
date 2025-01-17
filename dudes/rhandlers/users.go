package rhandlers

import (
	"fmt"
	"keepitok/dudes/models"
	"keepitok/dudes/utils"
	"log"
	"net/http"
	"os"
	"text/template"
)

func RootHandler(w http.ResponseWriter, r *http.Request) {

	t, err := template.ParseFiles(
		"templates/header.html",
		"templates/index.html",
		"templates/footer.html")

	if err != nil {
		fmt.Fprintf(os.Stderr, "Fprintf: %v\n", err)
	}
	t.ExecuteTemplate(w, "index", nil)
}

func RegistrationHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles(
		"templates/header.html",
		"templates/registration.html",
		"templates/footer.html")

	if err != nil {
		fmt.Fprintf(os.Stderr, "Fprintf: %v\n", err)
	}
	t.ExecuteTemplate(w, "register", nil)
}

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), 405)
		return
	}

	for key, val := range r.Form {
		fmt.Printf("%s = %s", key, val)

	}

	firstname := r.FormValue("first_name")
	lastname := r.FormValue("last_name")
	email := r.FormValue("email")
	password1 := r.FormValue("password1")
	password2 := r.FormValue("password2")

	if email == "" || firstname == "" || lastname == "" || password2 == "" || password1 != password2 {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	// Создаём юзера
	usr := models.CreateUser(firstname, lastname, email, password1)

	newURL := fmt.Sprintf("users/profile?id=%d", usr.ID)
	http.Redirect(w, r, newURL, http.StatusSeeOther)
}

func UsersListHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), 405)
		return
	}

	t, err := template.ParseFiles(
		"templates/header.html",
		"templates/users.html",
		"templates/footer.html")

	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	users, err := models.UsersList()

	if err != nil {
		http.Error(w, http.StatusText(523), 523)
		return
	}

	sessID := utils.Cookie(r, "sessID")
	var currentUser *models.User
	if sessID != "" {
		currentUser, err = models.GetCurrentUser(sessID)
		if err != nil {
			log.Println(err)
		}
	}

	data := map[string]interface{}{
		"users":       users,
		"currentUser": currentUser,
	}

	t.ExecuteTemplate(w, "users", data)
}

func ShowUserProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), 405)
		return
	}

	t, err := template.ParseFiles(
		"templates/header.html",
		"templates/profile.html",
		"templates/footer.html")

	userID := r.FormValue("id")

	usr, err := models.UserProfile(userID)

	if err != nil {
		log.Fatal(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	data := map[string]interface{}{
		"id":        usr.ID,
		"email":     usr.Email,
		"firstname": usr.Firstname,
		"lastname":  usr.Lastname,
	}

	t.ExecuteTemplate(w, "profile", data)
}

func RenderLoginTemplate(w http.ResponseWriter, data map[string]interface{}) {
	t, err := template.ParseFiles(
		"templates/header.html",
		"templates/login.html",
		"templates/footer.html")

	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	t.ExecuteTemplate(w, "login", data)
}

func LoginUserHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		RenderLoginTemplate(w, nil)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")
	log.Println(password)
	session, err := models.LoginUser(email, password)

	if err != nil || email == "" || password == "" {
		log.Println(err)
		data := map[string]interface{}{"message": "НЕ-А"}

		RenderLoginTemplate(w, data)
		return
	} else {

		newURL := fmt.Sprintf("users/profile?id=%d", session.UserID)

		sessCook := &http.Cookie{Name: "sessID",
			Value:    session.SessID,
			HttpOnly: false}

		http.SetCookie(w, sessCook)
		http.Redirect(w, r, newURL, http.StatusSeeOther)
	}
}

func SearchUserHandler(w http.ResponseWriter, r *http.Request) {

	log.Println(r.Method)
	if r.Method != "POST" {
		RenderLoginTemplate(w, nil)
		return
	}

	t, err := template.ParseFiles(
		"templates/header.html",
		"templates/users.html",
		"templates/footer.html")

	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	q := r.FormValue("search")
	users, err := models.SearchUsers(q)

	if err != nil {
		http.Error(w, http.StatusText(523), 523)
		return
	}

	// sessID := utils.Cookie(r, "sessID")
	sessID := "la-la-wrong-sessID"

	var currentUser *models.User
	if sessID != "" {
		currentUser, err = models.GetCurrentUser(sessID)
		if err != nil {
			log.Println(err)
		}
	}

	data := map[string]interface{}{
		"users":       users,
		"currentUser": currentUser,
	}

	t.ExecuteTemplate(w, "users", data)
}

func LogoutUserHandler(w http.ResponseWriter, r *http.Request) {

	sessID := utils.Cookie(r, "sessID")
	if sessID != "" {
		err := models.CloseSession(sessID)
		if err != nil {
			log.Println(err)
		}
	}
	sessCook := &http.Cookie{Name: "sessID", Value: "", HttpOnly: false}
	http.SetCookie(w, sessCook)
	http.Redirect(w, r, "login", http.StatusSeeOther)
}
