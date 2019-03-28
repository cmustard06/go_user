package views

import (
	"net/http"
	"go_user/sessions"
	"html/template"
	"fmt"
)

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		//检查是否有sessionid，没有就创建，有就更新一下访问时间，这样可以延长session的会话时间
		if s := sessions.SessionMgr.CheckSession(w, r); s == "" {
			//没有sessionid,开始一个新会话
			sessions.SessionMgr.StartSession(w, r)
		}
		t, _ := template.ParseFiles("template/html/login.html")
		t.Execute(w, nil)
		return
	} else if r.Method == "POST" {
		var sess string
		if sess = sessions.SessionMgr.CheckSession(w, r); sess == "" {
			////没有在内存中找到从cookie中获取的sessionid,开始一个新会话
			sess = sessions.SessionMgr.StartSession(w, r)
		}
		r.ParseForm()
		username := r.Form["username"][0]
		password := r.Form["password"][0]
		fmt.Println(username, password)
		if username == "" || password == "" {
			http.Redirect(w, r, "/login", 301)
			return
		}
		if username == "admin" && password == "123456" {

			if b := sessions.SessionMgr.SetSessionVal(sess, "username", username); b {
				http.Redirect(w, r, "/welcome", 301)
				return
			} else {
				http.Redirect(w, r, "/login", 301)
				return
			}
		} else {
			//登录信息错误
			http.Redirect(w, r, "/login", 301)
			return
		}
		return
	}
}

func Welcome(w http.ResponseWriter, r *http.Request) {
	var s string
	if s = sessions.SessionMgr.CheckSession(w, r); s == "" {
		//没有在内存中找到对应的sessionid,开始一个新会话
		sessions.SessionMgr.StartSession(w, r)
		//return
	}
	fmt.Printf("session:%s\n", s)
	testid, err := r.Cookie("TestId")
	if err != nil {
		http.Redirect(w, r, "/login", 301)
		return
	}
	val, b := sessions.SessionMgr.GetSessionVal(testid.Value, "username")
	if !b {
		http.Redirect(w, r, "/login", 301)
		return
	}
	fmt.Fprintf(w, "Welcome %s", val.(string))
	return
}
