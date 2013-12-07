package wmc

import (
	"net/http"
	"errors"
	"time"
	
	"appengine"
	"appengine/datastore"
)

type Like struct{
	FromID, ToID string
	Time time.Time
}

func likeHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	
	loginInfo := loginDetails(r)
	
	id := r.FormValue("id")
	if id == "" || loginInfo.User == nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	
	addLike(c, loginInfo.User.ID, id)
	
	http.Redirect(w, r, "/profile?id=" + id, http.StatusFound)
}


func addLike(c appengine.Context, fromID, toID string) {
	
	err := datastore.RunInTransaction(c, func(c appengine.Context) error {
		if alreadyLiked(c, fromID, toID) {
			return nil
		}
		
		p := retrieveProfile(c, toID)

		if p == nil || !p.Chef {
			panic(errors.New("No target profile"))
		}
		
		
		key := datastore.NewIncompleteKey(c, "Like", likeBookKey(c))
		
		_, err := datastore.Put(c, key, &Like{
			fromID,
			toID,
			time.Now(),
		})
		
		
		// TODO Shard counter
		p.Likes++
		
		check(err)
		
		updateProfile(c, toID, p)
		
		return nil
	}, &datastore.TransactionOptions{XG : true})
	
	check(err)
}

func alreadyLiked(c appengine.Context, fromID, toID string) bool {
	keys, err := datastore.NewQuery("Like").Ancestor(likeBookKey(c)).Filter("FromID=", fromID).Filter("ToID=", toID).KeysOnly().GetAll(c, nil)
	check(err)
	
	return len(keys) > 0
}

func likeBookKey(c appengine.Context) *datastore.Key {
	return datastore.NewKey(c, "LikeBook", "default_likebook", 0, nil)
}
