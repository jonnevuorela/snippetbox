package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

func home(writer http.ResponseWriter, request *http.Request) {
	if request.URL.Path != "/" {
		http.NotFound(writer, request)
		return
	}

	files := []string{
		"./ui/html/base.tmpl",
		"./ui/html/partials/nav.tmpl",
		"./ui/html/pages/home.tmpl",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Print(err.Error())
		http.Error(writer, "Internal Server Error", 500)
		return
	}

	err = ts.ExecuteTemplate(writer, "base", nil)
	if err != nil {
		log.Print(err.Error())
		http.Error(writer, "Internal Server Error", 500)
	}

	writer.Write([]byte("Hello from snippetbox"))
}

func snippetView(writer http.ResponseWriter, request *http.Request) {
	id, err := strconv.Atoi(request.URL.Query().Get("id"))
	if err != nil || id < 1 {
		http.NotFound(writer, request)
	}
	fmt.Fprintf(writer, "Display specific snippet with ID %d...", id)
}

func snippetCreate(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		writer.WriteHeader(405)
		http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	writer.Write([]byte("Create a new snippet..."))
}

func downloadHandler(writer http.ResponseWriter, request *http.Request) {
	http.ServeFile(writer, request, "./ui/static/file.zip")
}
