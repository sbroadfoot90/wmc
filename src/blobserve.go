package wmc

import (
	"appengine"
	"appengine/blobstore"
	"net/http"
)

func blobServeHandler(w http.ResponseWriter, r *http.Request) {
	blobstore.Send(w, appengine.BlobKey(r.FormValue("blobKey")))
}
