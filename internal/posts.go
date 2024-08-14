package internal

import (
	"fmt"
	"forum/forms"
	"forum/models"
	"net/http"
	"strconv"
)

func (app *application) createPostForm(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/post/create" {
		app.notFound(w, r)
		return
	}

	if r.Method == http.MethodPost {
		app.createPost(w, r)
		return
	}

	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		app.clientError(w, r, http.StatusMethodNotAllowed)
		return
	}

	c, err := app.categories.Latest()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.render(w, r, "create_page.html", &templateData{
		Form:       forms.New(nil),
		Categories: c,
	})
}

func (app *application) createPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, r, http.StatusBadRequest)
		return
	}
	form := forms.New(r.PostForm)
	form.Required("user_id", "title", "content", "categories")
	form.MaxLength("title", 100)

	c, err := app.categories.Latest()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	if !form.Valid() {
		app.render(w, r, "create_page.html", &templateData{
			Form:       form,
			Categories: c,
		})
		return
	}
	id, err := app.posts.Insert(form.Get("user_id"), form.Get("title"), form.Get("content"), "0", "0")
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	err = app.post_category.Insert(strconv.Itoa(id), r.PostForm["categories"])
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	sessionID, err := app.getSessionIDFromRequest(w, r)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	app.putSessionData(sessionID, "flash", "Post successfully created!")

	http.Redirect(w, r, fmt.Sprintf("/post?id=%d", id), http.StatusSeeOther)
}

func (app *application) post(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/post" {
		app.notFound(w, r)
		return
	}

	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		app.clientError(w, r, http.StatusMethodNotAllowed)
		return
	}

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w, r)
		return
	}

	s, err := app.posts.Get(int(id))
	if err == models.ErrNoRecord {
		app.notFound(w, r)
		return
	} else if err != nil {
		app.serverError(w, r, err)
		return
	}

	c, err := app.comments.Latest(id)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	categories, err := app.post_category.Get(id)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.render(w, r, "show_page.html", &templateData{
		Post:        s,
		Comments:    c,
		PCRelations: categories,
	})
}
