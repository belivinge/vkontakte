package internal

import (
	"net/http"

	"github.com/gofrs/uuid"
)

func (app *application) getSessionIDFromRequest(w http.ResponseWriter, r *http.Request) (string, error) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		if err == http.ErrNoCookie {
			newSessionID, createErr := app.createNewSession()
			if createErr != nil {
				return "", createErr
			}
			http.SetCookie(w, &http.Cookie{
				Name:  "session_id",
				Value: newSessionID,
				Path:  "/",
			})
			return newSessionID, nil
		}
		return "", err
	}
	return cookie.Value, nil
}

func (app *application) putSessionData(sessionID string, key string, value interface{}) {
	sessionData := app.getSession(sessionID)
	if sessionData == nil {
		sessionData = make(map[string]interface{})
	}
	sessionData[key] = value
	app.putSession(sessionID, sessionData)
}

func (app *application) createNewSession(userID ...int) (string, error) {
	sessionID, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	sessionIDStr := sessionID.String()

	app.sessionMutex.Lock()
	defer app.sessionMutex.Unlock()

	if app.sessionStore == nil {
		app.sessionStore = make(map[string]map[string]interface{})
	}

	if len(userID) > 0 {
		existingSessionID, exists := app.activeSessions[userID[0]]
		if exists {
			app.deleteSession(existingSessionID)
		}
	}

	app.sessionStore[sessionIDStr] = make(map[string]interface{})
	if app.activeSessions == nil {
		app.activeSessions = make(map[int]string)
	}
	if len(userID) > 0 {
		app.activeSessions[userID[0]] = sessionIDStr
	}

	return sessionIDStr, nil
}


func (app *application) getSession(sessionID string) map[string]interface{} {
	app.sessionMutex.Lock()
	defer app.sessionMutex.Unlock()
	sessionData, exists := app.sessionStore[sessionID]
	if !exists {
		return make(map[string]interface{})
	}
	return sessionData
}

func (app *application) putSession(sessionID string, data map[string]interface{}) {
	app.sessionMutex.Lock()
	defer app.sessionMutex.Unlock()
	app.sessionStore[sessionID] = data
}

func (app *application) deleteSession(sessionID string) {
	app.sessionMutex.Lock()
	defer app.sessionMutex.Unlock()
	delete(app.sessionStore, sessionID)
}

func (app *application) setSessionToken(sessionID, token string) {
	app.sessionMutex.Lock()
	defer app.sessionMutex.Unlock()

	session, exists := app.sessionStore[sessionID]
	if !exists {
		session = make(map[string]interface{})
		app.sessionStore[sessionID] = session
	}
	session["csrf_token"] = token
}

func (app *application) getSessionToken(sessionID string) (string, bool) {
	app.sessionMutex.Lock()
	defer app.sessionMutex.Unlock()

	sessionData, exists := app.sessionStore[sessionID]
	if !exists {
		return "", false
	}

	token, exists := sessionData["csrf_token"]
	if !exists {
		return "", false
	}

	csrfToken, ok := token.(string)
	if !ok {
		return "", false
	}

	return csrfToken, true
}

func (app *application) deleteSessionData(sessionID string, key string) {
	app.sessionMutex.Lock()
	defer app.sessionMutex.Unlock()

	sessionData, exists := app.sessionStore[sessionID]
	if !exists {
		return
	}
	delete(sessionData, key)
	if len(sessionData) == 0 {
		delete(app.sessionStore, sessionID)
	}
}

func (app *application) getSessionUserID(sessionID string) int {
	app.sessionMutex.Lock()
	defer app.sessionMutex.Unlock()

	sessionData, exists := app.sessionStore[sessionID]
	if !exists {
		return 0
	}

	userID, ok := sessionData["userID"].(int)
	if !ok {
		return 0
	}

	return userID
}
