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

type userLoginForm struct {
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

	_, err = a.users.Insert(userForm.Name, userForm.Email, userForm.Password)

	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			userForm.CheckField(false, "email", "Email is already in use. Please log-in.")
			data := a.newTemplateData(r)
			data.Form = userForm
			a.render(w, http.StatusUnprocessableEntity, "signup.tmpl", &data)
		} else {
			a.serverError(w, err)
		}
		return
	}

	a.sessionManager.Put(r.Context(), "flash", "Your signup was successful! Please log in.")
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
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
	data := a.newTemplateData(r)
	data.Form = userLoginForm{}

	if data.UserID > 0 {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	a.render(w, http.StatusOK, "login.tmpl", &data)
}

func (a *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	loginForm := userLoginForm{}
	err := a.decodePostForm(r, &loginForm)

	if err != nil {
		a.clientError(w, http.StatusBadRequest)
		return
	}

	validateUserLogin(&loginForm)

	if !loginForm.Valid() {
		data := a.newTemplateData(r)
		data.Form = loginForm
		a.render(w, http.StatusUnprocessableEntity, "login.tmpl", &data)
		return
	}

	userId, err := a.users.Authenticate(loginForm.Email, loginForm.Password)

	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			loginForm.AddNonFieldError("Email or Password is not correct.")
			data := a.newTemplateData(r)
			data.Form = loginForm
			a.render(w, http.StatusUnprocessableEntity, "login.tmpl", &data)
		} else {
			a.serverError(w, err)
		}
		return
	}

	err = a.sessionManager.RenewToken(r.Context())
	if err != nil {
		a.serverError(w, err)
		return
	}

	a.sessionManager.Put(r.Context(), "userID", userId)
	a.sessionManager.Put(r.Context(), "flash", "Welcome back to SnippetBox!")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func validateUserLogin(loginForm *userLoginForm) {
	loginForm.CheckField(validator.NotBlank(loginForm.Email), "email", models.ValidationMessageNotBlank)
	loginForm.CheckField(validator.MatchesRegex(loginForm.Email, validator.EmailRX), "email", "This field must contain a valid email address")
	loginForm.CheckField(validator.NotBlank(loginForm.Password), "password", models.ValidationMessageNotBlank)
}

func (a *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Logout a user")
}
