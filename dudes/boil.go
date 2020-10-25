package main

import (
	"fmt"
	"keepitok/dudes/rhandlers"
	"net/http"
)

func main() {
	fmt.Println("listening on port :3000")

	// Раздаём статику из папки assets

	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	// Регистрируем хэндлеры роутов

	http.HandleFunc("/", rhandlers.RootHandler)
	http.HandleFunc("/users", rhandlers.UsersListHandler)
	http.HandleFunc("/users/profile", rhandlers.ShowUserProfile)
	http.HandleFunc("/registration", rhandlers.RegistrationHandler)
	http.HandleFunc("/register-user", rhandlers.CreateUserHandler)
	http.HandleFunc("/login", rhandlers.LoginUserHandler)
	http.HandleFunc("/logout", rhandlers.LogoutUserHandler)
	http.HandleFunc("/search-user", rhandlers.SearchUserHandler)

	// Ждем запросы и раздаем

	http.ListenAndServe(":3000", nil)
}
