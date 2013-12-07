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
	c := appengine.NewContext(r)
	
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	
	loginInfo := loginDetails(r)
	p, id := targetUser(r)
	
	if u.User == nil || id == "" || p == nil || !p.Chef {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	
	if r.FormValue("like") == "yes" {
		AddLike(c, loginInfo.User.ID, id)
	}
	
	comment := Comment{
		r.FormValue("comment"),
		loginInfo.User.ID,
		id,
		time.Now(),
	}
	
	
	toKey := datastore.NewKey(c, "Profile", id, 0, nil)
	key := datastore.NewIncompleteKey(c, "Comment", toKey)
	
	_, err := datastore.Put(c, key, &comment)
	
	check(err)
	
	http.Redirect(w, r, "/profile?id="+id, http.StatusFound)
}