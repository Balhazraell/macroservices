package main

import (
	"fmt"
	"html/template"
	"net/http"
)

const HOST = "localhost"
const PORT = 5432

func init() {
	fsJS := http.FileServer(http.Dir("../static/js/"))

	http.Handle("/js/", http.StripPrefix("/js", fsJS))

	http.HandleFunc("/", returnIndex)
}

func returnIndex(response http.ResponseWriter, request *http.Request) {
	t, err := template.ParseFiles("../static/index.html")

	if err != nil {
		fmt.Fprintf(response, err.Error())
	}

	templateErr := t.ExecuteTemplate(response, "index", nil)

	if templateErr != nil {
		fmt.Fprintf(response, templateErr.Error())
		fmt.Fprintf(response, t.DefinedTemplates())
	}
}

func main() {
	http.ListenAndServe(":8081", nil)
	fmt.Println("Server is started...")
}
