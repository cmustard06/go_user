package main

import (
	"net/http"
	"go_user/views"
	"log"
)


func router(){
	//启动静态文件服务,需要删除这个前缀static，才能找到对应的文件
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/",views.Index)
	http.HandleFunc("/login",views.Login)
	http.HandleFunc("/welcome",views.Welcome)
}


func main(){
	//初始化sessions

	router()
	log.Println("starting listen...")
	log.Fatal(http.ListenAndServe(":8000",nil))
}
