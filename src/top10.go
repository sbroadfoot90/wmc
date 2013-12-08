package wmc

import (
	"appengine"
	"appengine/datastore"
	"net/http"
)

func topHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	n := 10
	q := datastore.NewQuery("Profile").Filter("Chef=", true).Order("-Likes").Limit(n)
	q2 := datastore.NewQuery("Profile").Filter("Chef=", true).Order("-Comments").Limit(n)
	loginInfo := loginDetails(r)

	profilesLikes := make([]Profile, 0, n)
	profilesComments := make([]Profile, 0, n)

	keysLikes, err := q.GetAll(c, &profilesLikes)

	check(err)

	keysComments, err := q2.GetAll(c, &profilesComments)

	check(err)

	outputToJsonOrTemplate(w, r, struct {
		LoginInfo           *LoginInfo
		ProfilesLikes       []Profile
		ProfilesComments    []Profile
		ProfileKeysLikes    []*datastore.Key
		ProfileKeysComments []*datastore.Key
	}{
		loginInfo,
		profilesLikes,
		profilesComments,
		keysLikes,
		keysComments,
	}, "top10")

}
