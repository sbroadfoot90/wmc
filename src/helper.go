package wmc

import (
	"appengine"
	"appengine/datastore"
	"appengine/user"

	"encoding/json"
	"net/http"
)

func outputToJsonOrTemplate(w http.ResponseWriter, r *http.Request, data interface{}, templateName string) {
	if r.FormValue("json") == "true" {
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		templates[templateName].ExecuteTemplate(w, "root", data)
	}
}

// check aborts the current execution if err is non-nil.
func check(err error) {
	if err != nil {
		panic(err)
	}
}

// errorHandler wraps the argument handler with an error-catcher that
// returns a 500 HTTP error if the request fails (calls check with err non-nil).
func errorHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err, ok := recover().(error); ok {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				// w.WriteHeader(http.StatusInternalServerError)
				// templates.ExecuteTemplate(w, "error.html", err)
			}
		}()
		fn(w, r)
	}
}

// Returns profile and user that is logged in
// If first argument is nil, user has no profile.
// If second argument is nil, user is not logged in.
func loginDetails(r *http.Request) (*LoginInfo) {
	c := appengine.NewContext(r)
	u := user.Current(c)
	
	
	if u == nil {
		url, err := user.LoginURL(c, "/") // TODO first time login url
		check(err)
		
		return &LoginInfo{nil, nil, url}
	}
	
	url, err := user.LogoutURL(c, "/")
	check(err)
	
	id := u.ID

	key := datastore.NewKey(c, "Profile", id, 0, nil)
	var p Profile

	err = datastore.Get(c, key, &p)
	// TODO Handle ErrNoSuchEntity
	if err == datastore.ErrNoSuchEntity {
		return &LoginInfo{nil, u, url}
	}
	check(err)

	return &LoginInfo{&p, u, url}
}

// Returns target user, and their id. If first argument is nil, user could not be found.
func targetUser(r *http.Request) (*Profile, string) {
	c := appengine.NewContext(r)
	id := r.FormValue("id")

	key := datastore.NewKey(c, "Profile", id, 0, nil)
	var p Profile

	err := datastore.Get(c, key, &p)
	if err == datastore.ErrNoSuchEntity {
		return nil, id
	}
	check(err)

	return &p, id
}

func Username(c appengine.Context, id string) string {
	key := datastore.NewKey(c, "Profile", id, 0, nil)
	var p Profile

	err := datastore.Get(c, key, &p)
	
	if err == datastore.ErrNoSuchEntity {
		return ""
	}
	
	check(err)
	
	return p.Name
}
