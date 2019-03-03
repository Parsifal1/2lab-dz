package main

import (
	"fmt"           // пакет для форматированного ввода вывода
	"html/template" // пакет для логирования
	"net/http"      // пакет для поддержки HTTP протокола
	// пакет для работы с  UTF-8 строками
)

var text string

func indexHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./index.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	fmt.Println(text)

	t.ExecuteTemplate(w, "index", text)

}

func drawHandler(w http.ResponseWriter, r *http.Request) {
	content := r.FormValue("content")

	text = text + content

	http.Redirect(w, r, "/", 302)
}

func main() {

	text = ""

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/draw", drawHandler)

	http.ListenAndServe(":3000", nil)
}
