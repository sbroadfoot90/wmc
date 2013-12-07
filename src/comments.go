package wmc

import (
	"net/http"
	"appengine"
	"appengine/datastore"
	"time"
)

type Comment struct{
	Comment string
	FromID, ToID string
	Time time.Time
}

func commentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusUnauthorized)
		return
	}
	
	loginInfo := loginDetails(r)
	p, id := targetUser(r)
	
	if p == nil || !p.Chef {
		http.Redirect(w, r, "/", http.StatusPreconditionFailed)
		return
	}
	
	comment := Comment{
		r.FormValue("comment"),
		loginInfo.User.ID,
		id,
		time.Now(),
	}
	c := appengine.NewContext(r)
	
	toKey := datastore.NewKey(c, "Profile", id, 0, nil)
	key := datastore.NewIncompleteKey(c, "Comment", toKey)
	
	_, err := datastore.Put(c, key, &comment)
	
	check(err)
	
	http.Redirect(w, r, "/profile?id="+id, http.StatusFound)
}