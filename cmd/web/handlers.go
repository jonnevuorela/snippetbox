package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

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
	data := app.newTemplateData(r)

	app.render(writer, http.StatusOK, "create.tmpl", data)
}

func (app *application) snippetCreatePost(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	if err != nil {
		app.clientError(writer, http.StatusBadRequest)
		return
	}

	title := request.PostForm.Get("title")
	content := request.PostForm.Get("content")

	expires, err := strconv.Atoi(request.PostForm.Get("expires"))
	if err != nil {
		app.clientError(writer, http.StatusBadRequest)
		return
	}

	fieldErrors := make(map[string]string)

	if strings.TrimSpace(title) == "" {
		fieldErrors["title"] = "This field cannot be blank"
	} else if utf8.RuneCountInString(title) > 100 {
		fieldErrors["title"] = "This field cannot be more than 100 characters long"
	}

	if strings.TrimSpace(content) == "" {
		fieldErrors["content"] = "This field cannot be blank"
	}

	if expires != 1 && expires != 7 && expires != 365 {
		fieldErrors["expire"] = "This field must equal 1, 7 or 365"
	}

	if len(fieldErrors) > 0 {
		fmt.Fprint(writer, fieldErrors)
		return
	}

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
