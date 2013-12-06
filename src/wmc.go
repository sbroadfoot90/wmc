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
		fmt.Fprint(w, loginURL)
		t.ExecuteTemplate(w, "root.tmpl", loginURL)
	} else {
		t.ExecuteTemplate(w, "rootloggedin.tmpl", u.String())	
	}
	

}

var t = template.Must(template.New("").ParseFiles("tmpl/root.tmpl", "tmpl/rootloggedin.tmpl"))