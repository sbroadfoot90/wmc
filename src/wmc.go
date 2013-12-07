package wmc

import (
	"html/template"
	"net/http"

	"appengine"
	"appengine/user"
)

func init() {
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
		
		t.ExecuteTemplate(w, "root.tmpl", loginURL)
	} else {
		logoutURL, err := user.LogoutURL(c, "/")

		check(err)

		t.ExecuteTemplate(w, "rootloggedin.tmpl", struct {
			User, LogoutURL string
		}{
			u.String(),
			logoutURL,
		})
	}

}

var t = template.Must(template.New("").ParseFiles("tmpl/root.tmpl", "tmpl/rootloggedin.tmpl", "tmpl/profile.tmpl", "tmpl/edit.tmpl"))
