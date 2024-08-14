package internal

import (
	"fmt"
	"forum/forms"
	"net/http"
	"strconv"
)

var (
	likesCount    int
	dislikesCount int
	postNum       int = 20
)

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w, r)
		return
	}

	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		app.clientError(w, r, http.StatusMethodNotAllowed)
		return
	}

	page_id := r.URL.Query().Get("page")

	if len(page_id) == 0 {
		http.Redirect(w, r, "/?page=1", http.StatusSeeOther)
		return
	}

	id, err := strconv.Atoi(page_id)
	if err != nil || id < 1 {
		app.notFound(w, r)
		return
	}

	s, err := app.posts.Latest()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	c, err := app.categories.Latest()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.render(w, r, "home_page.html", &templateData{
		Form:       forms.New(nil),
		Posts:      s,
		Categories: c,
	})
}

func (app *application) results(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/results" {
		app.notFound(w, r)
		return
	}

	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		app.clientError(w, r, http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		app.clientError(w, r, http.StatusBadRequest)
		return
	}
	form := forms.New(r.Form)
	form.RequiredAtLeastOne("categories", "created", "liked")

	if !form.Valid() {
		sessionID, err := app.getSessionIDFromRequest(w, r)
		if err != nil {
			app.serverError(w, r, err)
			return
		}
		app.putSessionData(sessionID, "flash", "Please select at least one filter.")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	page_id := r.URL.Query().Get("page")

	if len(page_id) == 0 {
		url := fmt.Sprintf("%s&page=1", r.URL.RequestURI())
		http.Redirect(w, r, url, http.StatusSeeOther)
		return
	}

	id, err := strconv.Atoi(page_id)
	if err != nil || id < 1 {
		app.notFound(w, r)
		return
	}

	s, err := app.posts.Filter(r.Form, app.post_reactions.FilterByLiked, app.post_category.FilterByCategories)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	c, err := app.categories.Latest()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	requestURI := r.URL.RequestURI()
	requestURI = requestURI[:len(requestURI)-1]

	app.render(w, r, "home_page.html", &templateData{
		Posts:      s,
		Categories: c,
		RequestURI: requestURI,
	})
}
