package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/julienschmidt/httprouter"
	"snippetbox.hafiz.com.ng/internal/models"
)

func (a *application) home(w http.ResponseWriter, r *http.Request) {
	latest, err := a.snippets.Latest()
	if err != nil {
		a.serverError(w, err)
		return
	}

	data := a.newTemplateData(r)
	data.LatestSnippets = latest
	a.render(w, http.StatusOK, "home.tmpl", &data)
}

func (a *application) snippetView(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		a.notFound(w)
		return
	}

	s, err := a.snippets.Get(id)

	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			a.notFound(w)
			return
		}
		a.serverError(w, err)
		return
	}

	data := a.newTemplateData(r)
	data.Snippet = s
	a.render(w, http.StatusOK, "view.tmpl", &data)

}

func (a *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	data := a.newTemplateData(r)
	data.Form = snippetCreateForm{
		Expires: 365,
	}
	a.render(w, http.StatusOK, "create.tmpl", &data)
}

type snippetCreateForm struct {
	Title       string
	Content     string
	Expires     int
	FieldErrors map[string]string
}

func (a *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		a.clientError(w, http.StatusBadRequest)
		return
	}

	expires, err := strconv.Atoi(r.PostForm.Get("expires"))
	if err != nil {
		a.clientError(w, http.StatusBadRequest)
		return
	}
	form := snippetCreateForm{
		Title:   r.PostForm.Get("title"),
		Content: r.PostForm.Get("content"),
		Expires: expires,
	}
	form.FieldErrors = validateFormFields(form)

	if len(form.FieldErrors) > 0 {
		data := a.newTemplateData(r)
		data.Form = form
		a.render(w, http.StatusUnprocessableEntity, "create.tmpl", &data)
		return
	}

	id, err := a.snippets.Insert(form.Title, form.Content, expires)
	if err != nil {
		a.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}

func validateFormFields(form snippetCreateForm) map[string]string {
	fieldErrors := make(map[string]string)
	if strings.TrimSpace(form.Title) == "" {
		fieldErrors["title"] = "This field cannot be blank"

	} else if utf8.RuneCountInString(form.Title) > 100 {
		fieldErrors["title"] = "This field cannot be more than 100 characters long"
	}

	if strings.TrimSpace(form.Content) == "" {
		fieldErrors["content"] = "This field cannot be blank"
	}

	if form.Expires != 1 && form.Expires != 7 && form.Expires != 365 {
		fieldErrors["expires"] = "This field must equal 1, 7, 365"
	}
	return fieldErrors
}
