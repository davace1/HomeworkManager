/*
作者：cz
时间：2020.8.5
内容：测试sql.go和studentWork.go
 */
package sql

import (
	"testing"
)

/*
func TestStudentWork_ShowPath(t *testing.T) {
	a:=StudentWork{}
	str:=a.ShowPath(2)
	if str!="\\homeworkFile\\1\\a.txt"{
		t.Fatal("Got" ,str," Want", "\\homeworkFile\\1\\a.txt")
	}
}
*/

func TestInit(t *testing.T){
	err:= db.Ping()
	if err!=nil{
		t.Fatal("got:",err,"want: nil")
	}
}

//必须放在最后
func TestClose(t *testing.T){
	err:=close()
	if err!=nil{
		t.Fatal("got:",err,"want: nil")
	}
}

