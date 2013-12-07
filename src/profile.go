package wmc

import (
	"net/http"
	"appengine"
)

type Foodie struct{
	Id string
	Name string
	Tagline string
}

func profileHandler(w http.ResponseWriter, r *http.Request) {
	f := Foodie{
		r.FormValue("id"),
		"Jacob",
		"live long and prosper",
	}
	t.ExecuteTemplate(w, "profile.tmpl", f)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	
	_ = c
}