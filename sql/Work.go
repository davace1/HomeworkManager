/*
作者：cz
时间：20208.6
内容：work的功能实现
创建（插入） 从sql中获取 修改结束时间  删除压缩
*/
package sql

import (
	"archive/zip"
	"homeworkManager/MyLog"
	"homeworkManager/config"
	"io"
	"os"
	"strconv"
	"time"
)

var teacher = config.Config.Teacher
//作业的实例化变量，用于存储作业的基本信息和链接
type Work struct {
	teachaerID int       //发布作业的教师id
	WorkID     int       //作业id 		插入数据库时补全，不用自行填写
	Name       string    //作业名称
	EndTime    time.Time //作业截止时间
	BeginTime  time.Time //作业开始时间	插入数据库时补全，不用自行填写
	Text       string    //备注
	IsEnd	   bool		//是否到达截止时间，若到达，则可用于下载 会在从数据库获取时赋值
	UpdateURL  string	//用于上传的链接，会在从数据库获取信息时设置
}

//创建一个Work 将对象保存至数据库
func SetWork(name string, endtime time.Time, text string) error {
	//插入数据
	_, err := db.Exec("INSERT INTO work VALUE(NULL,?,?,?,NOW(),?);",
		teacher.TeacherID, name, endtime, text)
	if err != nil {
		MyLog.WarnLog.Println("插入数据失败", err)
	}
	//获取id
	ans,err:=db.Query("SELECT LAST_INSERT_ID()")
	if err!=nil{
		MyLog.WarnLog.Println("查询出错",err)
		return nil
	}
	var workid int
	ans.Next()
	ans.Scan(&workid)
	err=ans.Err()
	if err!=nil{
		MyLog.WarnLog.Println("获取数据出错",err)
		return nil
	}

	str := "./homeworkFile/" + strconv.Itoa(workid)
	//创建文件夹
	err = os.Mkdir(str, os.ModePerm)
	if err != nil {
		str = "." + str
		err = os.Mkdir(str, os.ModePerm)
		if err != nil {
			MyLog.WarnLog.Println("创建文件夹失败 ", err)
		}
	}
	//记录
	return err
}

//根据ID获取Work
//从数据库中获取
func GetWork(workID int) (Work, error) {
	//从数据库中查询
	ans, err := db.Query("SELECT teacherid,name,endtime,begintime,text FROM work WHERE workID=?;", workID)
	defer ans.Close()
	if err != nil {
		MyLog.WarnLog.Println("获取work失败", err)
		return Work{}, err
	}
	//实例化
	a := Work{
		teachaerID: 0,
		WorkID:     workID,
		Name:       "",
		EndTime:    time.Time{},
		BeginTime:  time.Time{},
		Text:       "",
	}
	//从查询结果中获取数据
	ans.Next()
	err = ans.Scan(&a.teachaerID, &a.Name, &a.EndTime, &a.BeginTime, &a.Text)
	if err != nil {
		MyLog.WarnLog.Println("从记录中读取失败", err)
		return a, err
	}
	//设置结束判断
	if a.EndTime.Before(time.Now()){
		a.IsEnd=true
	}else {
		a.IsEnd=false
	}
	//设置上传链接
	id:=strconv.Itoa(a.WorkID)
	a.UpdateURL=config.Config.Server.IP+"/update?workid="+id
	return a, nil

}
//返回所有作业信息
func GetAllWork(teacherid int) ([]Work,error){
	//从数据库查询
	ans,err:=db.Query("SELECT workid,name,endtime,begintime,text FROM work WHERE teacherid=? ORDER BY begintime DESC;",teacherid)
	if err!=nil{
		MyLog.WarnLog.Println("查询出错",err)
		return nil, err
	}
	//作业切片
	var arr []Work
	//获取查询结果中所有数据
	for ans.Next(){
		var workid int
		var name,text string
		var endtime,begintime time.Time
		//输入到变量
		ans.Scan(&workid,&name,&endtime,&begintime,&text)
		//添加到切片
		arr=append(arr,Work{
			teachaerID: teacherid,
			WorkID:     workid,
			Name:       name,
			EndTime:    endtime,
			BeginTime:  begintime,
			Text:       text,
		})
	}
	//获取遇到的错误
	err=ans.Err()
	if err!=nil{
		MyLog.WarnLog.Println("获取作业失败")
		return nil, err
	}
	return arr,nil
}
//判断作业是否存在
func judgeWork(workID int) bool {
	//查询workID
	ans, err := db.Query("SELECT COUNT(*) 数量 FROM work where workID=?;", workID)
	defer ans.Close()
	if err != nil {
		MyLog.WarnLog.Println("查询作业号", workID, "失败:", err)
		return false
	}
	var num int
	//从查询结果中获取数据
	ans.Next()
	err = ans.Scan(&num)
	if err != nil {
		MyLog.WarnLog.Println("向数据库查询错误:", err)
		return false
	}
	if num != 0 {
		return true
	} else {
		return false
	}
}

//获取作业的名称
func getWorkName(workid int) (string, error) {
	ans, err := db.Query("SELECT name FROM work WHERE workID=?;", workid)
	defer ans.Close()
	if err != nil {
		MyLog.WarnLog.Println("查询数据库失败", err)
		return "", err
	}
	var name string
	//获取数据
	ans.Next()
	err = ans.Scan(&name)
	if err != nil {
		MyLog.WarnLog.Println("获取失败")
		return "", err
	}
	return name, nil
}

//以下为work成员函数



//用于删除
func (w Work) Delete() error {
	_, err := db.Exec("DELETE FROM work WHERE workid=?;", w.WorkID)
	if err != nil {
		MyLog.ErrorLog.Println("删除作业错误", err)
		return err
	}
	str := "./homeworkFile/" + w.Name
	//查看文件信息 用于判断文件存在
	_, err = os.Stat(str)
	if err != nil {
		str = "." + str
		_, err = os.Stat(str)
		if err != nil {
			MyLog.WarnLog.Println("文件不存在", err)
		}
	}
	//删除文件夹下所有文件
	err = os.RemoveAll(str)
	if err != nil {
		MyLog.WarnLog.Println("删除文件夹失败")
	}
	return nil
}

//获取已上传作业的全部路径
func (w Work) GetAllWorkPath() (allPath []string, err error) {
	//查询
	ans, err := db.Query("SELECT filepath FROM studentwork WHERE workid=? ;", w.WorkID)
	defer ans.Close()
	if err != nil {
		MyLog.WarnLog.Println("查询文件出错", err)
		return
	}
	//保存到切片
	for ans.Next() {
		var str string
		ans.Scan(&str)
		allPath = append(allPath, str)
	}
	err = ans.Err()
	if err != nil {
		MyLog.WarnLog.Println("查询文件出错", err)
		return
	}
	return
}

//将上传的作业文件压缩
func (w Work) ZIP() (string,string, error) {
	zipPath := "./zip/" + w.Name + ".zip"
	//创建压缩文件
	zipfile, err := os.Create(zipPath)
	if err != nil {
		zipPath="."+zipPath
		zipfile, err = os.Create(zipPath)
		if err!=nil{
			MyLog.WarnLog.Println("创建压缩文件出错", err)
			return "","", err
		}
	}
	defer zipfile.Close()
	//创建写入器
	archer := zip.NewWriter(zipfile)
	defer archer.Close()
	//获取改作业的所有文件的路径
	AllPath, err := w.GetAllWorkPath()
	if err != nil {
		return "","", err
	}
	//遍历所有文件
	for _, filename := range AllPath {
		file, err := os.Open(filename)
		if err != nil {
			filename = "." + filename
			file, err = os.Open(filename)
			if err != nil {
				MyLog.WarnLog.Println("打开文件失败", err)
				return "","", err
			}
		}
		defer  file.Close()
		first, err := os.Stat(filename)
		if err != nil {
			MyLog.WarnLog.Println("获取文件信息错误", err)
			return "","", err
		}
		//写入文件信息
		fh, err := zip.FileInfoHeader(first)
		if err != nil {
			MyLog.WarnLog.Println("不能", err)
			return "","", err
		}
		w, err := archer.CreateHeader(fh)
		if err != nil {
			MyLog.WarnLog.Println("不能创建添加器", err)
			return "","", err
		}
		//保存文件内容
		_, err = io.Copy(w, file)
		if err != nil {
			MyLog.WarnLog.Println("复制文件出错", err)
			return "","", err
		}
	}
	return zipPath,w.Name+".zip", nil
}
//获取提交作业的数量
func (w Work)GetUpdateNum()(int){
	//查询
	ans,err:=db.Query("SELECT COUNT(*) FROM studentwork WHERE workid=?;",w.WorkID)
	if err!=nil{
		MyLog.WarnLog.Println("查询数据库出错",err)
		return 0
	}
	var num int
	//获取查询结果
	ans.Next()
	ans.Scan(&num)
	err=ans.Err()
	if err!=nil{
		MyLog.WarnLog.Println("获取数据出错",err)
		return 0
	}
	return num
}
//修改作业的信息
func (w Work)Change(name string,endtime time.Time,text string)error{
	//修改
	_,err:=db.Exec("UPDATE work SET name=?,endtime=?,text=? WHERE WorkID=?;",
		name,endtime,text,w.WorkID)
	if err!=nil{
		MyLog.WarnLog.Println("修改数据失败",err)
		return err
	}
	return nil
}