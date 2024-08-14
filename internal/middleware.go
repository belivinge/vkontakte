package internal

import (
	"context"
	"fmt"
	"forum/models"
	"net/http"
)


func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.serverError(w, r, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("X-Frame-Options", "deny")

		next.ServeHTTP(w, r)
	})
}

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())

		next.ServeHTTP(w, r)
	})
}

func (app *application) requireAuthenticatedUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if app.authenticatedUser(r) == nil {
			http.Redirect(w, r, "user/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionID, err := app.getSessionIDFromRequest(w, r)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		sessionData := app.getSession(sessionID)
		userID, exists := sessionData["userID"]
		if !exists {
			next.ServeHTTP(w, r)
			return
		}

		user, err := app.users.Get(userID.(int))
		if err == models.ErrNoRecord {
			app.deleteSession(sessionID)
			next.ServeHTTP(w, r)
			return
		} else if err != nil {
			app.serverError(w, r, err)
			return
		}

		ctx := context.WithValue(r.Context(), contextKeyUser, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *application) standardMiddleware(next http.Handler) http.Handler {
	return app.recoverPanic(app.logRequest(secureHeaders(next)))
}

func (app *application) dynamicMiddleware(next http.Handler) http.Handler {
	return app.authenticate(next)
}
