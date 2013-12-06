package wmc

import (
	"fmt"
	"net/http"
	
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
		
		fmt.Fprint(w, "<html><body><a href=" + loginURL + ">login</a></body></html>")
	} else {
		fmt.Fprint(w, "Hello, you are logged in as ", u.String())
	}
	
}