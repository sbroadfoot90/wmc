package wmc

import (
	"net/http"
	"appengine"
	"appengine/user"
)

func firstTimeHandler(w http.ResponseWriter, r *http.Request){
	loginInfo :=  loginDetails(r)
	
	if loginInfo.Profile == nil && loginInfo.User != nil {
		http.Redirect(w , r , "/edit" , http.StatusFound)
	} else {
		redir := r.FormValue("redirect")
		if redir == "" {
			redir = "/"
		}
		http.Redirect(w, r, redir, http.StatusFound)		
	}
	
	
}


func LoginURL(c appengine.Context, dest string) (string, error) {
	url, err := user.LoginURL(c, "/firsttime?redirect=" + dest)
	return url, err
}