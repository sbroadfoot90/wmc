package wmc

import (
	"html/template"
	"net/http"

	"appengine"
	"appengine/user"
)

func init() {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/profile/", profileHandler)
	http.HandleFunc("/edit/", editHandler)
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
		logoutURL, err := user.LogoutURL(c, "/")

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		t.ExecuteTemplate(w, "rootloggedin.tmpl", struct {
			User, LogoutURL string
		}{
			u.String(),
			logoutURL,
		})
	}

}

var t = template.Must(template.New("").ParseFiles("tmpl/root.tmpl", "tmpl/rootloggedin.tmpl", "tmpl/profile.tmpl", "tmpl/edit.tmpl"))
