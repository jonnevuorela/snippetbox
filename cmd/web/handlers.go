package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"snippetbox.jonnevuorela.com/internal/models"

	"github.com/julienschmidt/httprouter"
)

func (app *application) home(writer http.ResponseWriter, request *http.Request) {
	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(writer, err)
		return
	}

	data := app.newTemplateData(request)
	data.Snippets = snippets

	app.render(writer, http.StatusOK, "home.tmpl", data)
}

func (app *application) snippetView(writer http.ResponseWriter, request *http.Request) {

	params := httprouter.ParamsFromContext(request.Context())

	id, err := strconv.Atoi(params.ByName("id"))
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

	data := app.newTemplateData(request)
	data.Snippet = snippet

	app.render(writer, http.StatusOK, "view.tmpl", data)

}

func (app *application) snippetCreate(writer http.ResponseWriter, r *http.Request) {
	writer.Write([]byte("Display the form for creating a new snippet"))
}

func (app *application) snippetCreatePost(writer http.ResponseWriter, request *http.Request) {
	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n- Kobayashi Issa"
	expires := 7

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(writer, err)
		return
	}

	http.Redirect(writer, request, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}

func downloadHandler(writer http.ResponseWriter, request *http.Request) {
	http.ServeFile(writer, request, "./ui/static/file.zip")
}
