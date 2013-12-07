package wmc

import (
	"net/http"
	"errors"
	
	"appengine"
)

func likeHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	
	loginInfo := loginDetails(r)
	
	id := r.FormValue("id")
	if id == "" || loginInfo.Profile == nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	
	
	AddLike(c, loginInfo.User.ID, id)
	
	http.Redirect(w, r, "/profile?id=" + id, http.StatusFound)
}


func AddLike(c appengine.Context, fromId, toID string) {
	p := retrieveProfile(c, toID)
	
	if p == nil || !p.Chef {
		panic(errors.New("No target profile"))
	}
	// TODO Transaction
	// TODO Shard counter
	p.Likes++
	
	updateProfile(c, toID, p)
}
