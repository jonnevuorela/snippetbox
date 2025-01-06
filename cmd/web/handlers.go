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

type snippetCreateForm struct {
	Title       string
	Content     string
	Expires     int
	FieldErrors map[string]string
}

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

	data.Form = snippetCreateForm{
		Expires: 365,
	}
	app.render(writer, http.StatusOK, "create.tmpl", data)
}

func (app *application) snippetCreatePost(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	if err != nil {
		app.clientError(writer, http.StatusBadRequest)
		return
	}

	expires, err := strconv.Atoi(request.PostForm.Get("expires"))
	if err != nil {
		app.clientError(writer, http.StatusBadRequest)
		return
	}

	form := snippetCreateForm{
		Title:       request.PostForm.Get("title"),
		Content:     request.PostForm.Get("content"),
		Expires:     expires,
		FieldErrors: map[string]string{},
	}

	if strings.TrimSpace(form.Title) == "" {
		form.FieldErrors["title"] = "This field cannot be blank"
	} else if utf8.RuneCountInString(form.Title) > 100 {
		form.FieldErrors["title"] = "This field cannot be more than 100 characters long"
	}

	if strings.TrimSpace(form.Content) == "" {
		form.FieldErrors["content"] = "This field cannot be blank"
	}

	if form.Expires != 1 && form.Expires != 7 && form.Expires != 365 {
		form.FieldErrors["expire"] = "This field must equal 1, 7 or 365"
	}

	if len(form.FieldErrors) > 0 {
		data := app.newTemplateData(request)
		data.Form = form
		app.render(writer, http.StatusUnprocessableEntity, "create.tmpl", data)
		return
	}

	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.serverError(writer, err)
		return
	}

	http.Redirect(writer, request, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}

func downloadHandler(writer http.ResponseWriter, request *http.Request) {
	http.ServeFile(writer, request, "./ui/static/file.zip")
}
