package internal

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func (app *application) handleDislike(w http.ResponseWriter, r *http.Request) {
	var request ReactionRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if r.URL.Path != request.URL {
		app.notFound(w, r)
		return
	}
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		app.clientError(w, r, http.StatusMethodNotAllowed)
		return
	}
	Num, err := strconv.Atoi(request.ID)
	count := 0
	var IfExist error
	if request.Type == "post" {
		count, IfExist = app.post_reactions.Insert(request.ID, request.UserID, "0")
		err = app.posts.UpdateReactions(Num, app.post_reactions.Likes, app.post_reactions.Dislikes)
	} else {
		count, IfExist = app.comment_reactions.Insert(request.ID, request.UserID, "0")
		err = app.comments.UpdateReactions(Num, app.comment_reactions.Likes, app.comment_reactions.Dislikes)
	}
	if err != nil || IfExist != nil {
		app.serverError(w, r, err)
		return
	}
	count1, _ := strconv.Atoi(request.Likes)
	count2, _ := strconv.Atoi(request.Dislikes)
	count2 = count2 + count
	if count == 0 {
		count2++
		count1--
	}
	response := struct {
		Likes    int `json:"likes"`
		Dislikes int `json:"dislikes"`
	}{
		Likes:    count1,
		Dislikes: count2,
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
}
