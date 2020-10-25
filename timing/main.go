package main

import (
	"fmt"
	"keepitok/timing/models"
	"net/http"
)

func main() {
	fmt.Println("Timing ms :: listening on port :3030")

	// Регистрируем хэндлеры роутов
	http.HandleFunc("/", models.GetUsersThingsHandler)

	// Регистрируем хэндлеры роутов
	http.HandleFunc("/add-thing", models.CreateThingHandler)
	// Инициализируем подключения к БД
	models.Init()

	// Ждем запросы и раздаем
	http.ListenAndServe(":3030", nil)
}
