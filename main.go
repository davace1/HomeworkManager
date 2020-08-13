package main

import (
	"homeworkManager/MyHandler"
	"homeworkManager/config"
	"net/http"
)

var  sc=config.Config.Server
func main(){
	//配置服务器的信息
	server:=&http.Server{
		Addr:              sc.Addr,
		Handler:           MyHandler.SetMux(),
	}
	//启动服务器
	server.ListenAndServe()
}