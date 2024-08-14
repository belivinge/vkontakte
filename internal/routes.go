package internal

import (
	"net/http"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/", app.dynamicMiddleware(http.HandlerFunc(app.home)))
	mux.Handle("/results", app.dynamicMiddleware(http.HandlerFunc(app.results)))
	mux.Handle("/post/create", app.dynamicMiddleware(app.requireAuthenticatedUser(http.HandlerFunc(app.createPostForm))))
	mux.Handle("/like/", app.dynamicMiddleware(app.requireAuthenticatedUser(http.HandlerFunc(app.handleLike))))
	mux.Handle("/dislike/", app.dynamicMiddleware(app.requireAuthenticatedUser(http.HandlerFunc(app.handleDislike))))
	mux.Handle("/comment/like/", app.dynamicMiddleware(app.requireAuthenticatedUser(http.HandlerFunc(app.handleLike))))
	mux.Handle("/comment/dislike/", app.dynamicMiddleware(app.requireAuthenticatedUser(http.HandlerFunc(app.handleDislike))))
	mux.Handle("/post/create/comment/", app.dynamicMiddleware(app.requireAuthenticatedUser(http.HandlerFunc(app.createComment))))
	mux.Handle("/post", app.dynamicMiddleware(http.HandlerFunc(app.post)))
	mux.Handle("/user/signup", app.dynamicMiddleware(http.HandlerFunc(app.signupUserForm)))
	mux.Handle("/user/login", app.dynamicMiddleware(http.HandlerFunc(app.loginUserForm)))
	mux.Handle("/user/logout", app.dynamicMiddleware(app.requireAuthenticatedUser(http.HandlerFunc(app.logoutUser))))
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))
	return app.standardMiddleware(mux)
}
