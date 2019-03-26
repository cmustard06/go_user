package views

import (
	"net/http"
	"go_user/sessions"
)

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method=="GET"{


	}else if r.Method == "POST"{
		if s := sessions.SessionMgr.CheckSession(w, r);s==""{
			http.Redirect(w,r,"/",301)
			return
		}
	}



}
