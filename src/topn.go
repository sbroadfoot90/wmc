package wmc

import (
	"appengine"
	"appengine/datastore"
	"net/http"
)

func topnHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	n := 10
	q := datastore.NewQuery("Profile").Filter("Chef=", true).Order("-Likes").Limit(n)
	loginInfo := loginDetails(r)

	profiles := make([]Profile, 0, n)

	keys, err := q.GetAll(c, &profiles)

	check(err)

	outputToJsonOrTemplate(w, r, struct {
		LoginInfo   *LoginInfo
		Profiles    []Profile
		ProfileKeys []*datastore.Key
	}{
		loginInfo,
		profiles,
		keys,
	}, "topn")

}
