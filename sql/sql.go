package sql

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"homeworkManager/MyLog"
	"homeworkManager/config"
	//手动加上mysql的库  go原生sql不支持mysql
)
//读取MySql配置信息
var dbc=config.Config.DataBase
//服务器
var db *sql.DB

//包级初始化，会在包初始化时自动调用
//用于创建与服务器的链接
func init(){
	DSN:=dbc.User+":"+dbc.Pw+"@("+dbc.Host+")/"+dbc.DBname+"?charset="+dbc.Charset+"&parseTime=true"
	var err error
	db,err=sql.Open(dbc.Dialect,DSN)
	if err!=nil{
		MyLog.ErrorLog.Println("数据库初始化失败:",err)
	}
}
//关闭数据库的链接
func close() error{
	err:=db.Close()
	if err!=nil{
		MyLog.ErrorLog.Println("关闭数据库链接失败:",err)
	}
	return err
}
