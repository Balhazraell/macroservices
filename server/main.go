package main

import (
	"fmt"
	"html/template"
	"net/http"

	"./core"
	"./logger"
	"./websockets"
)

func init() {
	fsJS := http.FileServer(http.Dir("../static/js/dist/"))
	fsCSS := http.FileServer(http.Dir("../static/css/"))

	http.Handle("/js/", http.StripPrefix("/js", fsJS))
	http.Handle("/css/", http.StripPrefix("/css", fsCSS))

	http.HandleFunc("/", returnIndex)
}

func returnIndex(response http.ResponseWriter, request *http.Request) {
	t, err := template.ParseFiles("../static/index.html")

	if err != nil {
		fmt.Fprintf(response, err.Error())
	}

	templateErr := t.ExecuteTemplate(response, "index.html", nil)

	if templateErr != nil {
		fmt.Fprintf(response, templateErr.Error())
		fmt.Fprintf(response, t.DefinedTemplates())
	}
}

func main() {
	defer logger.InfoPrint("Сервер закончил работу.")
	// Запускаем логгер.
	ok := logger.InitLogger()

	if ok {
		fmt.Println("Logger - YES")
	} else {
		fmt.Println("Logger - NO")
	}

	// Запускаем игровой сервер.
	core.GameServerStart()

	// Стартуем сервер websocket.
	websockets.Start()

	// Стартуем сервер статики. Стартуем его последним.
	// ListenAndServe - ждем завершения, по этому код дальше не выполняется.
	logger.InfoPrint("Сервер запущен и готов к работе.")
	http.ListenAndServe(":8081", nil)
}
