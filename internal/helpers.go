package internal

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"forum/models"
	"net/http"
	"runtime/debug"
	"time"
)

func generateCSRFToken() (string, error) {
	tokenBytes := make([]byte, 32) // 32 bytes = 256 bits
	_, err := rand.Read(tokenBytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(tokenBytes), nil
}

func (app *application) addDefaultData(w http.ResponseWriter, td *templateData, r *http.Request) *templateData {
	if td == nil {
		td = &templateData{}
	}
	sessionID, err := app.getSessionIDFromRequest(w, r)
	if err == nil {
		csrfToken, _ := app.getSessionToken(sessionID)
		td.CSRFToken = csrfToken
	} else {
		td.CSRFToken = ""
	}

	td.AuthenticatedUser = app.authenticatedUser(r)
	td.CurrentYear = time.Now().Year()

	if err == nil {
		flash := app.getSession(sessionID)["flash"]
		if flash != nil {
			td.Flash = flash.(string)
			app.deleteSessionData(sessionID, "flash")
		} else {
			td.Flash = ""
		}
	}

	return td
}

func (app *application) render(w http.ResponseWriter, r *http.Request, name string, td *templateData) {
	ts, ok := app.templatecache[name]
	if !ok {
		app.serverError(w, r, fmt.Errorf("The template %s does not exist", name))
		return
	}
	buf := new(bytes.Buffer)
	err := ts.Execute(buf, app.addDefaultData(w, td, r))
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	_, err = buf.WriteTo(w)
	if err != nil {
		panic(err)
	}
}

func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)

	error := &models.Error{
		Code:    http.StatusInternalServerError,
		Message: http.StatusText(http.StatusInternalServerError),
	}

	app.render(w, r, "error_page.html", &templateData{
		Error: error,
	})
}

func (app *application) clientError(w http.ResponseWriter, r *http.Request, status int) {
	err := &models.Error{
		Code:    status,
		Message: http.StatusText(status),
	}

	app.render(w, r, "error_page.html", &templateData{
		Error: err,
	})
}

func (app *application) notFound(w http.ResponseWriter, r *http.Request) {
	app.clientError(w, r, http.StatusNotFound)
}

func (app *application) authenticatedUser(r *http.Request) *models.User {
	value := r.Context().Value(contextKeyUser)

	user, ok := value.(*models.User)
	if !ok {
		return nil
	}

	return user
}
