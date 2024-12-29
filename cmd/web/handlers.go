package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"snippetbox.jonnevuorela.com/internal/models"
)

func (app *application) home(writer http.ResponseWriter, request *http.Request) {
	if request.URL.Path != "/" {
		app.notFound(writer)
		return
	}

	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(writer, err)
		return
	}
	for _, snippet := range snippets {
		fmt.Fprintf(writer, "%+v\n", snippet)
	}

	//	files := []string{
	//		"./ui/html/base.tmpl",
	//		"./ui/html/partials/nav.tmpl",
	//		"./ui/html/pages/home.tmpl",
	//	}
	//
	// ts, err := template.ParseFiles(files...)
	//
	//	if err != nil {
	//		app.serverError(writer, err)
	//		return
	//	}
	//
	// err = ts.ExecuteTemplate(writer, "base", nil)
	//
	//	if err != nil {
	//		app.serverError(writer, err)
	//	}
}

func (app *application) snippetView(writer http.ResponseWriter, request *http.Request) {
	id, err := strconv.Atoi(request.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(writer)
		return
	}
	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(writer)
		} else {
			app.serverError(writer, err)
		}
		return
	}

	fmt.Fprintf(writer, "%+v", snippet)
}

func (app *application) snippetCreate(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		writer.Header().Set("Allow", http.MethodPost)
		app.clientError(writer, http.StatusMethodNotAllowed)
		return
	}

	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n- Kobayashi Issa"
	expires := 7

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(writer, err)
		return
	}

	http.Redirect(writer, request, fmt.Sprintf("/snippet/view?id=%d", id), http.StatusSeeOther)
}

func downloadHandler(writer http.ResponseWriter, request *http.Request) {
	http.ServeFile(writer, request, "./ui/static/file.zip")
}
