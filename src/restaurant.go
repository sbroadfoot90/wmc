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
		return
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

func newRestaurantHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	loginInfo := loginDetails(r)

	if loginInfo.User == nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	} else {
		if r.Method == "GET" {
			editRestaurantGetHandler(w, loginInfo, nil, "")
		}

	}

}

func editRestaurantHandler(w http.ResponseWriter, r *http.Request) {
	loginInfo := loginDetails(r)

	if r.FormValue("rid") == "" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	rest, rid := targetRestaurant(r)

	if rid == "" {
		http.Error(w, "Restaurant Not Found", http.StatusNotFound)
		return
	}

	if loginInfo.User == nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	} else {
		if r.Method == "GET" {
			editRestaurantGetHandler(w, loginInfo, rest, rid)
		} else if r.Method == "POST" {
			editRestaurantPostHandler(w, r, loginInfo, rest, rid)
		}
	}
}

func editRestaurantGetHandler(w http.ResponseWriter, loginInfo *LoginInfo, rest *Restaurant, rid string) {
	templates["editRestaurant"].ExecuteTemplate(w, "root", struct {
		LoginInfo  *LoginInfo
		Restaurant *Restaurant
		RID        string
	}{
		loginInfo,
		rest,
		rid,
	})
}

// Handles creation and updating of restaurants
func editRestaurantPostHandler(w http.ResponseWriter, r *http.Request, loginInfo *LoginInfo, rest *Restaurant, rid string) {
	c := appengine.NewContext(r)
	if rest == nil {
		rest = &Restaurant{}
	}

	rest.Name = r.FormValue("Name")
	rest.Address = r.FormValue("Address")

	key := datastore.NewKey(c, "Restaurant", rid, 0, nil)

	_, err := datastore.Put(c, key, rest)

	check(err)

	http.Redirect(w, r, "/restaurant?rid="+rid, http.StatusFound)
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
