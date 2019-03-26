package main

import (
	"net/http"
	"go_user/views"
	"log"
)


func router(){
	http.HandleFunc("/",views.Index)
	http.HandleFunc("/login",views.Login)
}


func main(){
	//初始化sessions

	router()
	log.Println("starting listen...")
	log.Fatal(http.ListenAndServe(":8000",nil))
}
