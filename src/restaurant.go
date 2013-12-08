package wmc

import (
	"errors"
	"net/http"
	"strings"

	"appengine"
	"appengine/blobstore"
	"appengine/datastore"
	"appengine/memcache"
)

type Restaurant struct {
	Name           string
	Address        string
	URL            string
	RestaurantLogo appengine.BlobKey
}

func restaurantHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
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

	q := datastore.NewQuery("Profile").Filter("CurrentRestaurantID=", rid)
	profiles := make([]*Profile, 0, 100)

	keys, err := q.GetAll(c, &profiles)
	check(err)

	outputToJsonOrTemplate(w, r, struct {
		LoginInfo   *LoginInfo
		RID         string
		Restaurant  *Restaurant
		ProfileKeys []*datastore.Key
		Profiles    []*Profile
	}{
		loginInfo,
		rid,
		rest,
		keys,
		profiles,
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
			editRestaurantGetHandler(w, r, loginInfo)
		}
	}
}

func editRestaurantHandler(w http.ResponseWriter, r *http.Request) {
	loginInfo := loginDetails(r)
	if loginInfo.User == nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	} else {
		if r.Method == "GET" {
			editRestaurantGetHandler(w, r, loginInfo)
		} else if r.Method == "POST" {
			editRestaurantPostHandler(w, r, loginInfo)
		}
	}
}

func editRestaurantGetHandler(w http.ResponseWriter, r *http.Request, loginInfo *LoginInfo) {
	c := appengine.NewContext(r)
	uploadURL, err := blobstore.UploadURL(c, "/editRestaurant", nil)
	check(err)

	rest, rid := targetRestaurant(r)

	templates["editRestaurant"].ExecuteTemplate(w, "root", struct {
		LoginInfo  *LoginInfo
		Restaurant *Restaurant
		RID        string
		UploadURL  string
	}{
		loginInfo,
		rest,
		rid,
		uploadURL.String(),
	})
}

// Handles creation and updating of restaurants
func editRestaurantPostHandler(w http.ResponseWriter, r *http.Request, loginInfo *LoginInfo) {
	c := appengine.NewContext(r)

	blobs, values, err := blobstore.ParseUpload(r)
	check(err)
	rid := values.Get("rid")

	if rid == "" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	rest := retrieveRestaurant(c, rid)

	if rest == nil {
		rest = &Restaurant{}
	}

	rest.Name = values.Get("Name")
	rest.Address = values.Get("Address")
	rest.URL = values.Get("URL")

	var oldRestaurantLogo appengine.BlobKey

	if len(blobs["RestaurantLogo"]) > 0 {
		blobInfo := blobs["RestaurantLogo"][0]
		if !strings.HasPrefix(blobInfo.ContentType, "image") {
			blobstore.Delete(c, blobInfo.BlobKey) // discard error, doesn't matter as this is just an optimisation
			check(errors.New("File uploaded not image type"))
		}
		oldRestaurantLogo = rest.RestaurantLogo
		rest.RestaurantLogo = blobInfo.BlobKey
	}

	updateRestaurant(c, rid, rest)

	//if the update was successful, delete old restaurant logo
	if oldRestaurantLogo != "" {
		blobstore.Delete(c, oldRestaurantLogo) // discard error, doesn't matter as this is just an optimisation
	}

	http.Redirect(w, r, "/restaurant?rid="+rid, http.StatusFound)
}

func targetRestaurant(r *http.Request) (*Restaurant, string) {
	c := appengine.NewContext(r)
	rid := r.FormValue("rid")

	if rid == "" {
		return nil, ""
	}
	rest := retrieveRestaurant(c, rid)
	return rest, rid
}

func retrieveRestaurant(c appengine.Context, rid string) *Restaurant {
	if rid == "" {
		return nil
	}
	key := datastore.NewKey(c, "Restaurant", rid, 0, nil)
	var rest Restaurant

	_, err := memcache.Gob.Get(c, "restaurant-"+rid, &rest)

	if err == memcache.ErrCacheMiss {
		c.Debugf("Memcache Miss")
		err := datastore.Get(c, key, &rest)
		if err == datastore.ErrNoSuchEntity {
			return nil
		}
		check(err)
		memcache.Gob.Set(c, &memcache.Item{Key: "restaurant-" + rid, Object: rest})
	} else {
		check(err)
	}

	return &rest
}

func updateRestaurant(c appengine.Context, rid string, rest *Restaurant) {
	key := datastore.NewKey(c, "Restaurant", rid, 0, nil)
	_, err := datastore.Put(c, key, rest)
	memcache.Delete(c, "restaurant-"+rid)
	check(err)
}
