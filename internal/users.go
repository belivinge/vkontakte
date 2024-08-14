package internal

import (
	"fmt"
	"forum/forms"
	"forum/models"
	"net/http"
)

func (app *application) signupUserForm(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		app.signupUser(w, r)
		return
	}

	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		app.clientError(w, r, http.StatusMethodNotAllowed)
		return
	}

	app.render(w, r, "signup_page.html", &templateData{
		Form: forms.New(nil),
	})
}

func (app *application) signupUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, r, http.StatusBadRequest)
		return
	}
	form := forms.New(r.PostForm)
	form.Required("name", "phone", "email", "password")
	form.MatchesPattern("email", forms.EmailRX)
	form.MinLength("password", 10)

	if !form.Valid() {
		app.render(w, r, "signup_page.html", &templateData{
			Form: form,
		})
		return
	}
	err = app.users.Insert(form.Get("name"), form.Get("phone"), form.Get("email"), form.Get("password"))
	if err == models.ErrDuplicateEmail {
		form.Errors.Add("email", "Address is already in use")
		app.render(w, r, "signup_page.html", &templateData{
			Form: form,
		})
		return
	} else if err != nil {
		fmt.Println(err)
		fmt.Printf("**%s**", err)
		app.serverError(w, r, err)
		return
	}

	sessionID, err := app.getSessionIDFromRequest(w, r)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	app.putSessionData(sessionID, "flash", "Your signup was successful. Please log in.")
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (app *application) loginUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, r, http.StatusBadRequest)
		return
	}
	form := forms.New(r.PostForm)
	id, err := app.users.Authenticate(form.Get("email"), form.Get("password"))
	if err == models.ErrInvalidCredentials {
		form.Errors.Add("generic", "Email or Password is incorrect")
		app.render(w, r, "login_page.html", &templateData{
			Form: form,
		})
		return
	} else if err != nil {
		app.serverError(w, r, err)
		return
	}
	sessionID, err := app.createNewSession(id)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	app.putSessionData(sessionID, "userID", id)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) loginUserForm(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/user/login" {
		app.notFound(w, r)
		return
	}
	if r.Method == http.MethodPost {
		app.loginUser(w, r)
		return
	}

	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		app.clientError(w, r, http.StatusMethodNotAllowed)
		return
	}
	app.render(w, r, "login_page.html", &templateData{
		Form: forms.New(nil),
	})
}

func (app *application) logoutUser(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/user/logout" {
		app.notFound(w, r)
		return
	}

	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		app.clientError(w, r, http.StatusMethodNotAllowed)
		return
	}
	sessionID, err := app.getSessionIDFromRequest(w, r)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	userID := app.getSessionUserID(sessionID)
	app.deleteSession(sessionID)
	delete(app.activeSessions, userID)
	app.putSessionData(sessionID, "flash", "You've been logged out successfully!")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
