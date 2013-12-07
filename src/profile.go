package wmc

import (
	"net/http"

	"appengine"
	"appengine/datastore"
	"appengine/memcache"
)

type Profile struct {
	Name          string
	Tagline       string
	Chef          bool
	Title         string
	Likes         int
	RestaurantIds []string
}

// TODO populate from file
var Titles = []string{
	"Executive Chef",
	"Sous Chef",
	"Chef de Partie",
	"Demi Chef",
	"Pastry Chef",
	"Chef",
}

func profileHandler(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("id") == "" {
		http.Error(w, "Profile Not Found", http.StatusNotFound)
		// TODO Maybe have a list of all profiles here?
	}

	loginInfo := loginDetails(r)

	p, id := targetUser(r)

	if p == nil {
		http.Error(w, "Profile Not Found", http.StatusNotFound)
		return
	}

	c := appengine.NewContext(r)
	key := datastore.NewKey(c, "Profile", id, 0, nil)

	n := 10

	q := datastore.NewQuery("Comment").Ancestor(key).Order("-Time").Limit(n)

	comments := make([]Comment, 0, n)

	_, err := q.GetAll(c, &comments)

	check(err)

	outputToJsonOrTemplate(w, r, struct {
		LoginInfo    *LoginInfo
		ID           string
		Profile      *Profile
		Comments     []Comment
		AlreadyLiked bool
		C            appengine.Context
	}{
		loginInfo,
		id,
		p,
		comments,
		loginInfo.User == nil || alreadyLiked(c, loginInfo.User.ID, id),
		c,
	}, "profile")
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	loginInfo := loginDetails(r)

	if loginInfo.User == nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	} else {
		if loginInfo.Profile == nil {
			loginInfo.Profile = &Profile{}
		}
		if r.Method == "GET" {
			editGetHandler(w, loginInfo)
		} else if r.Method == "POST" {
			editPostHandler(w, r, loginInfo)
		}

	}

}

func editGetHandler(w http.ResponseWriter, loginInfo *LoginInfo) {
	templates["edit"].ExecuteTemplate(w, "root", struct {
		LoginInfo   *LoginInfo
		ValidTitles []string
		UploadURL   string
	}{
		loginInfo,
		Titles,
		"/edit",
	})
}

func editPostHandler(w http.ResponseWriter, r *http.Request, loginInfo *LoginInfo) {
	c := appengine.NewContext(r)

	p := loginInfo.Profile

	if p == nil {
		p = &Profile{}
	}

	p.Name = r.FormValue("Name")
	if tagline := r.FormValue("Tagline"); len(tagline) <= 40 {
		p.Tagline = tagline
	}

	isChef := r.FormValue("IsChef") == "yes"
	p.Chef = isChef
	if isChef {
		for _, title := range Titles {
			if title == r.FormValue("Title") {
				p.Title = title
			}
		}
	}

	updateProfile(c, loginInfo.User.ID, p)

	if isChef {
		http.Redirect(w, r, "/profile?id="+loginInfo.User.ID, http.StatusFound)
	} else {
		http.Redirect(w, r, "/root", http.StatusFound)
	}
}

func retrieveProfile(c appengine.Context, id string) *Profile {
	key := datastore.NewKey(c, "Profile", id, 0, nil)
	var p Profile

	_, err := memcache.Gob.Get(c, "profile-"+id, &p)

	if err == memcache.ErrCacheMiss {
		c.Debugf("Memcache Miss")
		err := datastore.Get(c, key, &p)
		if err == datastore.ErrNoSuchEntity {
			return nil
		}
		check(err)
		memcache.Gob.Set(c, &memcache.Item{Key: "profile-" + id, Object: p})
	} else {
		check(err)
	}

	return &p
}

func updateProfile(c appengine.Context, id string, p *Profile) {
	key := datastore.NewKey(c, "Profile", id, 0, nil)
	_, err := datastore.Put(c, key, p)
	memcache.Delete(c, "profile-"+id)
	check(err)
}
