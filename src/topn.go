package wmc

import (
	"appengine"
	"appengine/datastore"
	"net/http"
)


func topnHandler (w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	n := 10
	q := datastore.NewQuery("Profile").Order("-Likes").Limit(n)
	loginInfo := loginDetails(r)


	profiles := make([]Profile, 0, n)

	_, err := q.GetAll(c, &profiles)

	check(err)
	_  = loginInfo
	// outputToJsonOrTemplate(w, r, struct {
	// 	LoginInfo    *LoginInfo
	// 	Profiles	[]Profile
	// 	C            appengine.Context
	// }{
	// 	loginInfo,
	// 	profiles,
	// 	c,
	// }, "topn")
	// 
	// 
	
	
}
