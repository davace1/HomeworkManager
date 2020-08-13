/*
作者：cz
时间2020.8.7.
内容 测试Work
注意 注意对数据库的清空
 */

package sql

import (
	"os"
	"testing"
	"time"
)

func TestSetWork(t *testing.T) {
	err:=SetWork("test",time.Now(),"这是备注")
	if err!=nil{
		panic(err)
	}
	ans,err:=db.Query("SELECT COUNT(*),workid FROM work where name=?;","test")
	if err!=nil{
		panic(err)
	}
	ans.Next()
	var num,id int
	ans.Scan(&num,&id)
	if num==0{
		t.Fatal("Want 1 got ",num)
	}
	a,err:=GetWork(id)
	if err!=nil{
		panic(err)
	}
	a.Delete()
}

func TestWork_Delete(t *testing.T) {
	err:=SetWork("test",time.Now(),"这是备注")
	if err!=nil{
		panic(err)
	}
	ans,err:=db.Query("SELECT COUNT(*),workid FROM work where name=?;","test")
	if err!=nil{
		panic(err)
	}
	ans.Next()
	var num,id int
	ans.Scan(&num,&id)
	if num==0{
		panic(err)
	}
	a,err:=GetWork(id)
	if err!=nil{
		panic(err)
	}
	err=a.Delete()
	if err!=nil{
		t.Fatal("Want nil got ",err)
	}
}


func TestWork_ZIP(t *testing.T) {
	err:=SetWork("test",time.Now(),"这是备注")
	if err!=nil{
		panic(err)
	}
	ans,err:=db.Query("SELECT COUNT(*),workid FROM work where name=?;","test")
	if err!=nil{
		panic(err)
	}
	ans.Next()
	var num,id int
	ans.Scan(&num,&id)
	if num==0{
		panic(err)
	}
	a,err:=GetWork(id)
	if err!=nil{
		panic(err)
	}
	SetStudentWork(a.workID,"hello.c",[]byte("hello world!"))
	SetStudentWork(a.workID,"world.c",[]byte("hello world!"))
	_,err=a.ZIP()
	if  err!=nil{
		t.Fatal("Want nil got ",err)
	}
	str:="../ZIP/"+a.Name+".zip"
	_,err=os.Stat(str)
	if err!=nil{
		t.Fatal("Zip文件不存在")
	}
	err=a.Delete()
	if err!=nil{
		t.Fatal("Want nil got ",err)
	}
}