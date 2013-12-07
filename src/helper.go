package wmc

import (
	"appengine"
	"appengine/datastore"
	
	"net/http"
	"encoding/json"
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

func loginDetails(r *http.Request) (*Foodie, bool){
	c:= appengine.NewContext(r)
	id := r.FormValue("id")
	c.Debugf(id)
	key := datastore.NewKey(c, "Profile", id, 0, nil)
	var f Foodie
	
	err := datastore.Get(c, key, &f)
	// TODO Handle ErrNoSuchEntity
	if err == datastore.ErrNoSuchEntity {
		return nil, false
	}
	check(err)
	
	return &f, true
}