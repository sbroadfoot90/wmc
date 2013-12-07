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
	user, _ := loginDetails(r)
	
	outputToJsonOrTemplate(w, r, user, "profile")
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	u := user.Current(c)

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
			_, err := datastore.Put(c, key, &foodie)
			
			check(err)
			
			http.Redirect(w, r, "/profile?id="+u.ID, http.StatusFound)
		}

	}
}
