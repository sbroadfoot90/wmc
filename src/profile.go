package wmc

import (
	"net/http"
	
	"appengine"
	"appengine/datastore"
	"appengine/user"
)

type Foodie struct {
	Name    string
	Tagline string
}


func profileHandler(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("id") == "" {
		http.Error(w, "Profile Not Found", http.StatusNotFound)
		// TODO Maybe have a list of all profiles here?
	}
	
	user, id := targetUser(r)
	
	if user == nil {
		http.Error(w, "Profile Not Found", http.StatusNotFound)
		return
	}
	
	outputToJsonOrTemplate(w, r, struct{
		ID string
		User *Foodie
	}{
		id,
		user,
	}, "profile")
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	u, id = loginDetails(r)

	if u == nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	} else {
		key := datastore.NewKey(c, "Profile", u.ID, 0, nil)
		var foodie Foodie
				
		if r.Method == "GET" {
			err := datastore.Get(c, key, &foodie)
			// TODO include handling for ErrNoSuchEntity
			
			if err != datastore.ErrNoSuchEntity {
				check(err)
			}
			
			templates["edit"].ExecuteTemplate(w, "root", struct{
				ID     string
				Foodie Foodie
			}{
				u.ID,
				foodie,
			})
		} else if r.Method == "POST" {
			foodie.Name = r.FormValue("Name")
			foodie.Tagline = r.FormValue("Tagline")
			
			c := appengine.NewContext(r)
			_, err := datastore.Put(c, key, &foodie)
			
			check(err)
			
			http.Redirect(w, r, "/profile?id="+u.ID, http.StatusFound)
		}

	}
}
