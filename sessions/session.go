package sessions

import (
	"sync"
	"time"
	"io"
	"crypto/rand"
	"strconv"
	"encoding/base64"
	"net/http"
	"net/url"
)

var SessionMgr *SessionManager = nil


/*Session 管理*/
type SessionManager struct {
	CookieName string
	Lock sync.Mutex
	MaxLifeTime int64
	Sessions map[string]*Session   //{sessionId:session}
}

type Session struct {
	SessionId string
	LastTimeAccessed time.Time   //记录该会话上一次访问的时间
	Values map[interface{}]interface{}
}

//程序开始时初始化session管理器
func NewSessionManager(cookieName string, maxLifeTime int64) *SessionManager{
	mgr := &SessionManager{
		CookieName: cookieName,
		MaxLifeTime:maxLifeTime,
		Sessions: make(map[string] *Session), //初始化字典
	}
	//在初始化的同时，需要启动一个协程监控session已经存活的时间
	go mgr.GC()

	return mgr
}

func (mgr *SessionManager)GC(){
	mgr.Lock.Lock()
	defer mgr.Lock.Unlock()

	for sesionId, session := range mgr.Sessions{
		//删除超过时间的session
		if session.LastTimeAccessed.Unix()+mgr.MaxLifeTime<time.Now().Unix(){
			delete(mgr.Sessions,sesionId)
		}
	}
	time.AfterFunc(time.Duration(mgr.MaxLifeTime)*time.Second,mgr.GC)
}

//创建sessionID
func (mgr *SessionManager)GenerateSessionId()string{
	b := make([]byte,64)
	if _,err := io.ReadFull(rand.Reader,b);err!=nil{
		nano := time.Now().UnixNano()
		return strconv.FormatInt(nano,10)
	}
	return base64.URLEncoding.EncodeToString(b)
}


//开始一个新的session,主要用于新的会话连接
func (mgr *SessionManager)StartSession(w http.ResponseWriter, r *http.Request) string {
	mgr.Lock.Lock()
	defer mgr.Lock.Unlock()

	newSessionId := url.QueryEscape(mgr.GenerateSessionId())

	session := &Session{
		SessionId:newSessionId,
		LastTimeAccessed:time.Now(),
		Values:make(map[interface{}]interface{}),
	}
	mgr.Sessions[newSessionId] = session

	//对客户端的cookie设置过期时间
	cookie := http.Cookie{
		Name:mgr.CookieName,
		Value:newSessionId,
		Path:"/",
		HttpOnly:true,
		MaxAge:int(mgr.MaxLifeTime),
	}
	http.SetCookie(w, &cookie)
	return newSessionId
}

//结束会话，包括浏览器
func (mgr *SessionManager)EndSession(w http.ResponseWriter,r *http.Request){
	cookie,err := r.Cookie(mgr.CookieName)
	if err!=nil || cookie.Value == ""{
		return
	}else{
		mgr.Lock.Lock()
		defer mgr.Lock.Unlock()

		delete(mgr.Sessions, cookie.Value)
		//最后同步告知浏览器，cookie无效
		expire := time.Now()
		cookie := &http.Cookie{
			Name:mgr.CookieName,
			Path:"/",
			HttpOnly:true,
			Expires:expire,
			MaxAge:-1,
		}
		http.SetCookie(w,cookie)
	}
}

//服务器端删除无效的sessionid
func (mgr *SessionManager)DeleteSession(sessionId string){
	mgr.Lock.Lock()
	defer mgr.Lock.Unlock()
	delete(mgr.Sessions,sessionId)
}

//获取服务器中session相关值
func (mgr *SessionManager)GetSessionVal(sessionId string, key interface{})(interface{}, bool){
	mgr.Lock.Lock()
	defer mgr.Lock.Unlock()

	if session,ok := mgr.Sessions[sessionId];ok{
		if val,ok:=session.Values[key];ok{
			return val,ok
		}
	}
	return nil,false
}

//设置session的值
func(mgr *SessionManager)SetSessionVal(sessionId string, key interface{}, value interface{})bool{
	mgr.Lock.Lock()
	defer mgr.Lock.Unlock()

	if session,ok := mgr.Sessions[sessionId];ok{
		session.Values[key] = value
		return true
	}
	return false
}

//获取所有的sessionid 列表
func (mgr *SessionManager)GetSessionIDList() []string{
	mgr.Lock.Lock()
	defer mgr.Lock.Unlock()

	sessionIdList := make([]string,0)
	for k,_ := range mgr.Sessions{
		sessionIdList = append(sessionIdList, k)
	}
	return sessionIdList
}

//session检测
func (mgr *SessionManager)CheckSession(w http.ResponseWriter,r *http.Request) string {
	var cookie,err = r.Cookie(mgr.CookieName)
	if cookie == nil||err!=nil{
		return ""
	}
	mgr.Lock.Lock()
	defer mgr.Lock.Unlock()

	sessionId := cookie.Value
	if session,ok :=mgr.Sessions[sessionId];ok{
		session.LastTimeAccessed = time.Now()
		return sessionId
	}
	return ""
}

//获取最新的访问session更新时间
func (mgr *SessionManager)GetLastAccessTime(sessionId string) time.Time{
	mgr.Lock.Lock()
	defer mgr.Lock.Unlock()

	if session,ok := mgr.Sessions[sessionId];ok{
		return session.LastTimeAccessed
	}
	return time.Now()
}

func init(){
	SessionMgr = NewSessionManager("TestId", 60)
}