package wmc

import (
	"net/http"

	"appengine"
	"appengine/user"
)

func init() {
	importTemplates("tmpl")
	
	http.HandleFunc("/", errorHandler(rootHandler))
	http.HandleFunc("/profile", errorHandler(profileHandler))
	http.HandleFunc("/edit", errorHandler(editHandler))
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	u := user.Current(c)

	if u == nil {

		loginURL, err := user.LoginURL(c, "/")

		check(err)
		
		templates["index"].ExecuteTemplate(w, "root", loginURL)
	} else {
		logoutURL, err := user.LogoutURL(c, "/")

		check(err)

		templates["indexLoggedIn"].ExecuteTemplate(w, "root", struct {
			User, LogoutURL string
		}{
			u.String(),
			logoutURL,
		})
	}

}