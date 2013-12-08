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

type Profile struct {
	Name                string
	Tagline             string
	ProfilePicture      appengine.BlobKey
	Chef                bool
	Title               string
	Likes               int
	CurrentRestaurantID string
	PastRestaurantIds   []string
	Comments            int
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
	c := appengine.NewContext(r)

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

	// When viewing a chef, show comments made on chef's page
	// When viewing a foodie, show comments made to chefs
	var filterString string
	if p.Chef {
		filterString = "ToID="
	} else {
		filterString = "FromID="
	}

	n := 10
	q := datastore.NewQuery("Comment").Ancestor(commentBookKey(c)).Order("-Time").Filter(filterString, id).Limit(n)
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

// Edit get requests show an edit form, whereas post requests update a user.
func editHandler(w http.ResponseWriter, r *http.Request) {
	loginInfo := loginDetails(r)

	// If not logged in, cannot edit their profile
	if loginInfo.User == nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	} else {
		if loginInfo.Profile == nil {
			loginInfo.Profile = &Profile{}
		}
		if r.Method == "GET" {
			editGetHandler(w, r, loginInfo)
		} else if r.Method == "POST" {
			editPostHandler(w, r, loginInfo)
		}
	}
}

func editGetHandler(w http.ResponseWriter, r *http.Request, loginInfo *LoginInfo) {
	c := appengine.NewContext(r)

	uploadURL, err := blobstore.UploadURL(c, "/edit", nil)
	check(err)

	restaurants := make([]Restaurant, 0, 20)

	q := datastore.NewQuery("Restaurant").Limit(20)

	keys, err := q.GetAll(c, &restaurants)

	templates["edit"].ExecuteTemplate(w, "root", struct {
		LoginInfo      *LoginInfo
		ValidTitles    []string
		Restaurants    []Restaurant
		RestaurantKeys []*datastore.Key
		UploadURL      string
	}{
		loginInfo,
		Titles,
		restaurants,
		keys,
		uploadURL.String(),
	})
}

func editPostHandler(w http.ResponseWriter, r *http.Request, loginInfo *LoginInfo) {
	c := appengine.NewContext(r)

	blobs, values, err := blobstore.ParseUpload(r)
	check(err)

	p := loginInfo.Profile

	if p == nil {
		p = &Profile{}
	}

	p.Name = values.Get("Name")
	if tagline := values.Get("Tagline"); len(tagline) <= 40 {
		p.Tagline = tagline
	}

	isChef := values.Get("IsChef") == "yes"
	p.Chef = isChef
	if isChef {
		// validate chef title
		for _, title := range Titles {
			if title == values.Get("Title") {
				p.Title = title
			}
		}

		// validate restaurant id
		rid := values.Get("Restaurant")
		if rid != "" && retrieveRestaurant(c, rid) != nil {
			if p.CurrentRestaurantID != "" {
				alreadyExists := false
				for _, prid := range p.PastRestaurantIds {
					if rid == prid {
						alreadyExists = true
						break
					}
				}
				if !alreadyExists {
					p.PastRestaurantIds = append(p.PastRestaurantIds, rid)
				}
			}
			p.CurrentRestaurantID = rid

		}
	}
	var oldProfilePicture appengine.BlobKey

	if len(blobs["ProfilePicture"]) > 0 {
		blobInfo := blobs["ProfilePicture"][0]
		if !strings.HasPrefix(blobInfo.ContentType, "image") {
			blobstore.Delete(c, blobInfo.BlobKey) // discard error, doesn't matter as this is just an optimisation
			check(errors.New("File uploaded not image type"))
		}
		oldProfilePicture = p.ProfilePicture
		p.ProfilePicture = blobInfo.BlobKey
	}

	updateProfile(c, loginInfo.User.ID, p)

	//if the update was successful, delete old profile picture
	if oldProfilePicture != "" {
		blobstore.Delete(c, oldProfilePicture) // discard error, doesn't matter as this is just an optimisation
	}

	if isChef {
		http.Redirect(w, r, "/profile?id="+loginInfo.User.ID, http.StatusFound)
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
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
