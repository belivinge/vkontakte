package internal

import (
	"fmt"
	"forum/forms"
	"net/http"
)

func (app *application) createComment(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/post/create/comment/" {
		app.notFound(w, r)
		return
	}

	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		app.clientError(w, r, http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		app.clientError(w, r, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("post_id", "user_id", "content")

	if !form.Valid() {
		sessionID, err := app.getSessionIDFromRequest(w, r)
		if err != nil {
			app.serverError(w, r, err)
			return
		}
		app.putSessionData(sessionID, "flash", "Please type something into the comment section.")
		http.Redirect(w, r, fmt.Sprintf("/post?id=%s", form.Get("post_id")), http.StatusSeeOther)
		return
	}

	err = app.comments.Insert(form.Get("post_id"), form.Get("user_id"), form.Get("content"))
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	sessionID, err := app.getSessionIDFromRequest(w, r)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	app.putSessionData(sessionID, "flash", "Comment successfully created!")

	http.Redirect(w, r, fmt.Sprintf("/post?id=%s", form.Get("post_id")), http.StatusSeeOther)
}
