package views

import (
	"net/http"
	"io"
	"log"
	"go_user/sessions"
)

func Index(w http.ResponseWriter,r *http.Request){
	if r.Method == "GET"{
		//先检查session是否有效，在做是否重新生成session的
		if ok := sessions.SessionMgr.CheckSession(w,r);ok==""{
			sessions.SessionMgr.StartSession(w,r)
		}
		log.Print(sessions.SessionMgr.GetSessionIDList())
		w.WriteHeader(http.StatusOK)
		io.WriteString(w,"hello world")
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
	return
}
