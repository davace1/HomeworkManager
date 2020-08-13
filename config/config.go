/*
作者：Davace
时间：2020.8.3
作用：从config.json中读取数据，存至Config变量中，用于各类数据的初始化
 */

package config

import (
	"encoding/json"
	"homeworkManager/MyLog"
	"io/ioutil"
	"os"
)
//配置星系
type sqlmessage struct {
	Dialect string	`json:"dialect"`
	User 	string	`json:"user"`
	Pw		string	`json:"pw"`
	Host	string	`json:"host"`
	DBname	string	`json:"dbname"`
	Charset	string	`json:"charset"`
}
type teacher struct {
	TeacherID		int	`json:"teacherid,string"`	//教师id
	LogInName		string	`json:"loginname"`//登录名
	LogInPW			string 	`json:"loginpw"`//登录密码
	LastName		string	`json:"lastname"`//姓氏
}
type serverMassage struct {
	Addr 	string `json:"addr"`
	IP		string `json:"ip"`
}
type config struct {
	DataBase	sqlmessage `json:"database"`
	Teacher		teacher		`json:"teacher"`
	Server		serverMassage`json:"server"`
}
//用于访问配置信息
var Config config
//包初始化函数 会被最先调用
func init(){
	//由于测试和正式运行的目录不同，所以要尝试在当前目录和上级目录打开
	jsonFile,err:=os.Open("../config.json")
	if err!=nil{
		jsonFile,err=os.Open("./config.json")
		if err!=nil{
			MyLog.ErrorLog.Fatal(err)

		}
	}
	defer jsonFile.Close()
	//读取数据（字节）
	jsonDate,err:=ioutil.ReadAll(jsonFile)
	if err!=nil{
		MyLog.ErrorLog.Fatal(err)
	}
	//json解析
	json.Unmarshal(jsonDate,&Config)
}
