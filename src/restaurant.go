package wmc

import (
	"net/http"

	"appengine"
	"appengine/datastore"
)

type Restaurant struct {
	Name    string
	Address string
}

func restaurantHandler(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("rid") == "" {
		http.Error(w, "Restaurant Not Found", http.StatusNotFound)
	}

	loginInfo := loginDetails(r)

	rest, rid := targetRestaurant(r)

	if rest == nil {
		http.Error(w, "Restaurant Not Found", http.StatusNotFound)
		return
	}

	outputToJsonOrTemplate(w, r, struct {
		LoginInfo  *LoginInfo
		RID        string
		Restaurant *Restaurant
	}{
		loginInfo,
		rid,
		rest,
	}, "restaurant")
}

func targetRestaurant(r *http.Request) (*Restaurant, string) {
	c := appengine.NewContext(r)
	rid := r.FormValue("rid")

	rest := retrieveRestaurant(c, rid)

	return rest, rid
}

func retrieveRestaurant(c appengine.Context, rid string) *Restaurant {
	key := datastore.NewKey(c, "Restaurant", rid, 0, nil)
	var rest Restaurant

	err := datastore.Get(c, key, &rest)

	if err == datastore.ErrNoSuchEntity {
		return nil
	}
	check(err)

	return &rest
}
