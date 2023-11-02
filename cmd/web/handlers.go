package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"snippetbox.hafiz.com.ng/internal/models"
	"snippetbox.hafiz.com.ng/internal/validator"
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
	Title               string `form:"title"`
	Content             string `form:"content"`
	Expires             int    `form:"expires"`
	validator.Validator `form:"-"`
}

func (a *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	var form snippetCreateForm

	err := a.decodePostForm(r, &form)
	if err != nil {
		a.clientError(w, http.StatusBadRequest)
		return
	}

	validateNewSnippet(&form)

	if !form.Valid() {
		data := a.newTemplateData(r)
		data.Form = form
		a.render(w, http.StatusUnprocessableEntity, "create.tmpl", &data)
		return
	}

	id, err := a.snippets.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		a.serverError(w, err)
		return
	}

	a.sessionManager.Put(r.Context(), "flash", "Snippet successfully created!")

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}

func validateNewSnippet(form *snippetCreateForm) {
	form.CheckField(validator.NotBlank(form.Title), "title", models.ValidationMessageNotBlank)
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
	form.CheckField(validator.NotBlank(form.Content), "content", models.ValidationMessageNotBlank)
	form.CheckField(validator.PermittedInt(form.Expires, 1, 7, 365), "expires", "This field must equal 1, 7, 365")
}

type userSignupForm struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

func (a *application) userSignup(w http.ResponseWriter, r *http.Request) {
	data := a.newTemplateData(r)
	data.Form = userSignupForm{}
	a.render(w, http.StatusOK, "signup.tmpl", &data)
}

func (a *application) userSignupPost(w http.ResponseWriter, r *http.Request) {
	userForm := userSignupForm{}
	err := a.decodePostForm(r, &userForm)
	if err != nil {
		a.clientError(w, http.StatusBadRequest)
		return
	}

	validateNewUser(&userForm)
	if !userForm.Valid() {
		data := a.newTemplateData(r)
		data.Form = userForm
		a.render(w, http.StatusUnprocessableEntity, "signup.tmpl", &data)
		return
	}

	fmt.Fprint(w, "Sign up the user with the details provided")
}

func validateNewUser(userForm *userSignupForm) {
	userForm.CheckField(validator.NotBlank(userForm.Name), "name", models.ValidationMessageNotBlank)
	userForm.CheckField(validator.MaxChars(userForm.Name, 100), "name", "This field cannot be more than 100 characters long")
	userForm.CheckField(validator.NotBlank(userForm.Email), "email", models.ValidationMessageNotBlank)
	userForm.CheckField(validator.MatchesRegex(userForm.Email, validator.EmailRX), "email", "This field must contain a valid email address")
	userForm.CheckField(validator.NotBlank(userForm.Password), "password", models.ValidationMessageNotBlank)
	userForm.CheckField(validator.MinChars(userForm.Password, 8), "password", "This field must be at least 8 characters long")
}

func (a *application) userLogin(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Display a HTML form for logging a user in")
}

func (a *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Authenticate the user")
}

func (a *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Logout a user")
}
