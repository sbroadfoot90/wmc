package wmc

import (
	"appengine"
	"appengine/datastore"
	"net/http"
	"time"
	"errors"
)

type Comment struct {
	Comment      string
	FromID, ToID string
	Time         time.Time
}

// Handles posting of comments along with "also like"
func commentHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	loginInfo := loginDetails(r)
	p, id := targetUser(r)

	if loginInfo.User == nil || id == "" || p == nil || !p.Chef {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	if r.FormValue("like") == "yes" {
		addLike(c, loginInfo.User.ID, id)
	}

	comment := Comment{
		r.FormValue("comment"),
		loginInfo.User.ID,
		id,
		time.Now(),
	}

	key := datastore.NewIncompleteKey(c, "Comment", commentBookKey(c))

	 _, err := datastore.Put(c, key, &comment)
	 
	 check(err)

	addComment(c, comment.Comment, loginInfo.User.ID, id)
	
	http.Redirect(w, r, "/profile?id="+id, http.StatusFound)
}

func addComment(c appengine.Context, comment, fromID, toID string, ) {

	err := datastore.RunInTransaction(c, func(c appengine.Context) error {
		
		p := retrieveProfile(c, toID)

		if p == nil || !p.Chef {
			panic(errors.New("No target profile"))
		}

		key := datastore.NewIncompleteKey(c, "Comment", commentBookKey(c))

		_, err := datastore.Put(c, key, &Comment{
			comment,
			fromID,
			toID,
			time.Now(),
		})

		// TODO Shard counter
		p.Comments++

		check(err)

		updateProfile(c, toID, p)

		return nil
	}, &datastore.TransactionOptions{XG: true})

	check(err)
}

func commentBookKey(c appengine.Context) *datastore.Key {
	return datastore.NewKey(c, "commentBook", "default_commentbook", 0, nil)
}
