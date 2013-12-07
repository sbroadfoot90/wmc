package wmc

import (
	"net/http"

	"appengine"
	"appengine/datastore"
)

type Profile struct {
	Name       string
	Tagline    string
	Chef       bool
	RestuarantIds []string
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

	outputToJsonOrTemplate(w, r, struct {
		ID   string
		User *Profile
	}{
		id,
		user,
	}, "profile")
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	u, id := loginDetails(r)


	if u == nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	} else {
		if r.Method == "GET" {
			editGetHandler(w, u, id)
		} else if r.Method == "POST" {
			editPostHandler(w, r, u, id)
		}

	}
}

func editGetHandler(w http.ResponseWriter, u *Profile, id string) {
	templates["edit"].ExecuteTemplate(w, "root", struct {
		Profile *Profile
		ID     string
	}{
		u,
		id,
	})
}

func editPostHandler(w http.ResponseWriter, r *http.Request, u *Profile, id string) {
	c := appengine.NewContext(r)
	u.Name = r.FormValue("Name")
	u.Tagline = r.FormValue("Tagline")
	
	key := datastore.NewKey(c, "Profile", id, 0, nil)
	c.Debugf(id)
	_, err := datastore.Put(c, key, u)

	check(err)

	http.Redirect(w, r, "/profile?id="+id, http.StatusFound)
}
