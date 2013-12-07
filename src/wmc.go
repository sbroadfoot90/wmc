package wmc

import (
	"net/http"
//	"appengine"
//	"appengine/user"
)

func init() {
	importTemplates("tmpl")

	http.HandleFunc("/", errorHandler(rootHandler))
	http.HandleFunc("/profile", errorHandler(profileHandler))
	http.HandleFunc("/edit", errorHandler(editHandler))
	http.HandleFunc("/comment", errorHandler(commentHandler))
	http.HandleFunc("/firsttime", errorHandler(firstTimeHandler))
	http.HandleFunc("/like", errorHandler(likeHandler))
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	
	loginInfo := loginDetails(r)
	
	
	
	if loginInfo.Profile != nil && loginInfo.Profile.Chef {
		http.Redirect(w, r, "/profile?id=" + loginInfo.User.ID, http.StatusFound)

	} else{
		templates["index"].ExecuteTemplate(w, "root", struct{
			LoginInfo *LoginInfo
		}{
			loginInfo,
		})
		
	}
	
	
}
