package wmc

import (
	"net/http"
	"html/template"
	"fmt"
	
	"appengine"
	"appengine/user"
)

func init() {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/profile", profileHandler)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	
	u := user.Current(c)

	if u == nil {

		loginURL, err := user.LoginURL(c, "/")
		
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		t.ExecuteTemplate(w, "root.tmpl", loginURL)
	} else {
		t.ExecuteTemplate(w, "rootloggedin.tmpl", u.String())	
	}
	

}

var t = template.Must(template.New("").ParseFiles("tmpl/root.tmpl", "tmpl/rootloggedin.tmpl", "tmpl/profile.tmpl"))