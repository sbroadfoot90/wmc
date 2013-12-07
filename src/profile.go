package wmc

import (
	"net/http"

	"appengine"
	"appengine/datastore"
)

type Profile struct {
	Name       	string
	Tagline    	string
	Chef       	bool
	Title		string
	RestaurantIds []string
}

func profileHandler(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("id") == "" {
		http.Error(w, "Profile Not Found", http.StatusNotFound)
		// TODO Maybe have a list of all profiles here?
	}
	
	loginInfo := loginDetails(r)

	u, id := targetUser(r)

	if u == nil {
		http.Error(w, "Profile Not Found", http.StatusNotFound)
		return
	}
	
	c := appengine.NewContext(r)
	key := datastore.NewKey(c, "Profile", id, 0, nil)
	
	n := 10
	
	q := datastore.NewQuery("Comment").Ancestor(key).Order("-Time").Limit(n)
	
	comments := make([]Comment, 0, n)
	
	_, err := q.GetAll(c, &comments)
	
	check(err)

	outputToJsonOrTemplate(w, r, struct {
		LoginInfo *LoginInfo
		ID   string
		User *Profile
		Comments []Comment
		C appengine.Context
	}{
		loginInfo,
		id,
		u,
		comments,
		c,
	}, "profile")
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	loginInfo := loginDetails(r)


	if loginInfo.User == nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	} else {
		if loginInfo.Profile == nil {
			loginInfo.Profile = &Profile{}
		}
		if r.Method == "GET" {
			editGetHandler(w, loginInfo)
		} else if r.Method == "POST" {
			editPostHandler(w, r, loginInfo)
		}

	}
}

func editGetHandler(w http.ResponseWriter, loginInfo *LoginInfo) {
	templates["edit"].ExecuteTemplate(w, "root", struct{LoginInfo *LoginInfo}{loginInfo})
}

func editPostHandler(w http.ResponseWriter, r *http.Request, loginInfo *LoginInfo) {
	c := appengine.NewContext(r)
	
	loginInfo.Profile.Name = r.FormValue("Name")
	loginInfo.Profile.Tagline = r.FormValue("Tagline")
	loginInfo.Profile.Chef = r.FormValue("IsChef") == "yes"
	
	key := datastore.NewKey(c, "Profile", loginInfo.User.ID, 0, nil)
	c.Debugf(loginInfo.User.ID)
	_, err := datastore.Put(c, key, loginInfo.Profile)

	check(err)

	http.Redirect(w, r, "/profile?id="+loginInfo.User.ID, http.StatusFound)
}
