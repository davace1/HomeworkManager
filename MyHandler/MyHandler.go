/*
2020.8.8
cz
处理器函数和处理器内容
*/
package MyHandler

import (
	"homeworkManager/MyLog"
	"homeworkManager/config"
	"homeworkManager/sql"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"time"
)
import "net/http"

var teacher = config.Config.Teacher

func SetMux() *http.ServeMux {
	MUX := http.NewServeMux()
	//绑定处理器
	MUX.HandleFunc("/", index)
	MUX.HandleFunc("/login",login)
	MUX.HandleFunc("/hello",HelloTeacher)
	MUX.HandleFunc("/newWork",SetNewWork)
	MUX.HandleFunc("/work",ShowWorkMassage)
	MUX.HandleFunc("/create",createWork)
	MUX.HandleFunc("/delete",deleteWork)
	MUX.HandleFunc("/download",download)
	MUX.HandleFunc("/change",change)
	MUX.HandleFunc("/acceptChange",acceptChange)
	MUX.HandleFunc("/update",update)
	MUX.HandleFunc("/acceptUpdate",acceptupdate)
	return MUX
}
//返回错误页面
func SetError(w http.ResponseWriter,massage string){
	t,_:=template.ParseFiles("./MyHandler/error.html")
	t.Execute(w,massage)
	MyLog.InfoLog.Println("返回错误信息页面，内容：",massage)
	return
}
//判断登录，登录返回true
func JudgeLogged(r *http.Request)bool{
	//获取cookie 如果cookie记录有登录信息，则直接跳转到hello
	cookie,err:=r.Cookie("logged")
	if err==http.ErrNoCookie{
		//没有找到cookie
		return false
	}
	if err!=nil{
		MyLog.WarnLog.Println("获取cookie出错",err)
		return false
	}
	str,err:=strconv.Atoi(cookie.Value)
	if err!=nil{
		MyLog.ErrorLog.Println("数值转换出错",err)
		return false
	}
	if str==teacher.TeacherID{
		return true
	}

	return false
}
//欢迎界面，用于登录
func index(w http.ResponseWriter, r *http.Request) {
	//检查cookie 是否登录
	if JudgeLogged(r){
		//跳转链接
		http.Redirect(w,r,"/hello",http.StatusFound)
		HelloTeacher(w,r)
		return
	}
	t, err := template.ParseFiles("./MyHandler/index.html","./MyHandler/index_empty.html")
	if err != nil {
		MyLog.ErrorLog.Println("获取模板出错",err)
		SetError(w,"内部错误")
		return
	}
	err=t.ExecuteTemplate(w,"index","")
	if err!=nil{
		MyLog.WarnLog.Println("配置模板信息出错",err)
		SetError(w,"内部错误")
		return
	}

}
//获取等于的信息，验证登录信息正确性，并发送cookie到用户
func login(w http.ResponseWriter, r *http.Request) {
	//从表单获取登录信息，表单格式为applecation/x-www-for,-urlencoded,账号密码分别叫user,password
	user := r.PostFormValue("user")
	password := r.PostFormValue("password")
	if !(user == teacher.LogInName && password == teacher.LogInPW) {
		//账号密码错误
		//将登录界面和登录错误的界面相结合，返还给客户端
		t, err := template.ParseFiles("./MyHandler/index.html","./MyHandler/index_error.html")
		if err != nil {
			MyLog.ErrorLog.Println("获取模板出错",err)
			SetError(w,"内部错误")
			return
		}
		err=t.ExecuteTemplate(w,"index","")
		if err!=nil{
			MyLog.WarnLog.Println("配置模板信息出错",err)
			SetError(w,"内部错误")
			return
		}
		MyLog.InfoLog.Println("登录失败，登录账户为:",user)
	} else {
		//登陆成功，设定cokie，并渲染新的页面
		cookie:=http.Cookie{
			Name:       "logged",
			Value:      strconv.Itoa(teacher.TeacherID),
			MaxAge:     60*60*24*30,//一个月
			HttpOnly:   false,
		}
		http.SetCookie(w,&cookie)
		//登录成功并跳转
		http.Redirect(w,r,"/hello",http.StatusFound)
		HelloTeacher(w,r)
		MyLog.InfoLog.Println("登陆成功，登录账号为:",user)
		/*
		t, err := template.ParseFiles("./MyHandler/index.html","./MyHandler/index_empty.html")
		if err != nil {
			MyLog.ErrorLog.Println("获取模板出错",err)
		}
		err=t.ExecuteTemplate(w,"index","")
		if err!=nil{
			MyLog.WarnLog.Println("配置模板信息出错",err)
		}*/
	}
}
//渲染教师管理界面
func HelloTeacher(w http.ResponseWriter, r *http.Request){
	if !JudgeLogged(r){
		//未登录则重定向到默认页面
		index(w,r)
		return
	}
	arr,err:=sql.GetAllWork(teacher.TeacherID)
	if err!=nil{
		return
	}
	t,err:=template.ParseFiles("./MyHandler/Hello.html")
	if err!=nil{
		MyLog.WarnLog.Println("打开模板错误",err)
		SetError(w,"内部错误")
		return
	}
	err=t.Execute(w,arr)
	if err!=nil{
		MyLog.WarnLog.Println("插入数据出错",err)
		SetError(w,"内部错误")
		return
	}
}
//返回新建作业的页面
func SetNewWork(w http.ResponseWriter,r*http.Request){
	//先检查登录
	if !JudgeLogged(r) {
		index(w,r)
		return
	}
	t,err:=template.ParseFiles("./MyHandler/NewWork.html")
	if err!=nil{
		MyLog.WarnLog.Println("打开模板错误",err)
		SetError(w,"内部错误")
		return
	}
	err=t.Execute(w,"")
	if err!=nil{
		MyLog.ErrorLog.Println("配置模板出错",err)
		SetError(w,"内部错误")
		return
	}
}
//显示对应work的具体信息和功能
func ShowWorkMassage(w http.ResponseWriter,r*http.Request){
	//检查登录 避免直接访问
	if !JudgeLogged(r){
		index(w,r)
		return
	}
	//获取模板
	t,err:=template.ParseFiles("./MyHandler/WorkMassage.html")
	if err!=nil{
		MyLog.WarnLog.Println("打开模板错误",err)
		SetError(w,"内部错误")
		return
	}
	//解析workid 从url中
	//获取数据
	num,err:=strconv.Atoi(r.FormValue("id"))
	if err!=nil{
		MyLog.WarnLog.Println("转换出错",err)
		SetError(w,"内部错误")
		return
	}
	//http.Redirect(w,r,"/work",http.StatusFound)
	work,err:=sql.GetWork(num)
	if err!=nil{
		SetError(w,"作业不存在")
		return
	}
	err=t.Execute(w,work)
	if err!=nil{
		MyLog.ErrorLog.Println("配置模板出错",err)
		SetError(w,"内部错误")
		return
	}
	MyLog.InfoLog.Println("查看作业信息，作业id",work.WorkID)
}
//从浏览器获取数据，创建新的作业
func createWork(w http.ResponseWriter,r*http.Request){
	//检查登录 避免直接访问
	if !JudgeLogged(r){
		index(w,r)
		return
	}
	var name,text,timeStr string
	var endtime time.Time
	name=r.FormValue("name")
	timeStr=r.FormValue("endtime")
	text=r.FormValue("text")
	endtime,err:=time.Parse("2006-01-02",timeStr)
	if err!=nil{
		MyLog.WarnLog.Println("格式转换出错")
		SetError(w,"内部错误")
		return
	}
	err=sql.SetWork(name,endtime,text)
	if err!=nil{
		return
	}
	http.Redirect(w,r,"/hello",http.StatusFound)
	MyLog.InfoLog.Println("创建新的作业，作业名称:",name)
}
//删除work
func deleteWork(w http.ResponseWriter,r*http.Request){
	//检查登录 避免直接访问
	if !JudgeLogged(r){
		index(w,r)
		return
	}
	//解析workid 从url中
	//获取数据
	num,err:=strconv.Atoi(r.FormValue("workid"))
	if err!=nil{
		MyLog.WarnLog.Println("转换出错",err)
		SetError(w,"内部错误")
		return
	}
	//获取结构体(由于需要删除文件，所以必须获取
	a,err:=sql.GetWork(num)
	if err!=nil{
		SetError(w,"请检查是否存在作业")
		return
	}
	//执行删除，会从sql和本地文件中删除
	err=a.Delete()
	if err!=nil{
		SetError(w,"内部错误")
		return
	}
	//重定向到前一页面
	http.Redirect(w,r,"/hello",http.StatusFound)
	MyLog.InfoLog.Println("删除作业，作业的ID：",num)
}
//下载文件
func download(w http.ResponseWriter,r*http.Request){
	//检查登录 避免直接访问
	if !JudgeLogged(r){
		index(w,r)
		return
	}
	//解析workid 从url中
	r.ParseForm()
	//获取数据
	str:=r.FormValue("workid")
	num,err:=strconv.Atoi(str)
	if err!=nil{
		MyLog.WarnLog.Println("转换出错",err)
		SetError(w,"内部错误")
		return
	}
	//获取结构体(由于需要删除文件，所以必须获取
	a,err:=sql.GetWork(num)
	if err!=nil{
		SetError(w,"请检查是否存在作业")
		return
	}
	//创建压缩文件
	str,name,err:=a.ZIP()
	if err!=nil{
		SetError(w,"内部错误")
		return
	}
	file,err:=os.Open(str)
	if err!=nil{
		MyLog.WarnLog.Println("打开文件失败",err)
		SetError(w,"内部错误")
		return
	}
	defer file.Close()
	//设置响应的header头
	w.Header().Add("Content-type","application/octet-steam")
	w.Header().Add("content-disposition","attachment; filename=\""+name+"\"")
	_,err=io.Copy(w,file)
	if err!=nil{
		MyLog.WarnLog.Println("复制出错，上传出错",err)
		SetError(w,"内部错误")
		return
	}
	http.Redirect(w,r,"/hello",http.StatusFound)
	MyLog.InfoLog.Println("下载作业的文件，作业ID：",num)
}
//修改信息
func change(w http.ResponseWriter,r*http.Request){
	//检查登录 避免直接访问
	if !JudgeLogged(r){
		index(w,r)
		return
	}
	t,err:=template.ParseFiles("./MyHandler/ChangeWork.html")
	if err!=nil{
		MyLog.WarnLog.Println("打开模板错误",err)
		SetError(w,"内部错误")
		return
	}
	//解析workid 从url中
	err=r.ParseForm()
	if err!=nil{
		MyLog.WarnLog.Println("解析出错",err)
		SetError(w,"内部错误")
		return
	}
	//获取数据
	num,err:=strconv.Atoi(r.FormValue("workid"))
	if err!=nil{
		MyLog.WarnLog.Println("转换出错",err)
		SetError(w,"内部错误")
		return
	}
	//获取结构体(由于需要删除文件，所以必须获取
	a,err:=sql.GetWork(num)
	if err!=nil{
		return
	}
	err=t.Execute(w,a)
	if err!=nil{
		MyLog.ErrorLog.Println("配置模板出错",err)
		SetError(w,"内部错误")
		return
	}

}
//接受修改的信息，保存并更改
func acceptChange(w http.ResponseWriter,r*http.Request){
	//检查登录 避免直接访问
	if !JudgeLogged(r){
		index(w,r)
		return
	}
	r.ParseForm()
	var name,text,timeStr string
	var endtime time.Time
	name=r.FormValue("name")
	timeStr=r.FormValue("endtime")
	text=r.FormValue("text")
	endtime,err:=time.Parse("2006-01-02",timeStr)
	if err!=nil{
		MyLog.WarnLog.Println("格式转换出错")
		SetError(w,"内部错误")
		return
	}
	num,err:=strconv.Atoi(r.FormValue("workid"))
	if err!=nil{
		MyLog.WarnLog.Println("转换出错",err)
		SetError(w,"内部错误")
		return
	}
	a,err:=sql.GetWork(num)
	if err!=nil{
		SetError(w,"请检查作业存在")
		return
	}
	err=a.Change(name,endtime,text)
	if err!=nil{
		return
	}
	http.Redirect(w,r,"/hello",http.StatusFound)
	MyLog.InfoLog.Println("修改作业的信息,id:",num )
}
//上传作业页面
func update(w http.ResponseWriter,r*http.Request){
	//获取id
	r.ParseForm()
	num,err:=strconv.Atoi(r.FormValue("workid"))
	if err!=nil{
		MyLog.WarnLog.Println("转换出错",err)
		return
	}
	a,err:=sql.GetWork(num)
	if err!=nil{
		SetError(w,"请检查作业是否存在")
		return
	}
	t,err:=template.ParseFiles("./MyHandler/update.html")
	if err!=nil{
		MyLog.WarnLog.Println("模板打开失败",err)
		SetError(w,"内部错误")
		return
	}
	err=t.Execute(w,a)
	if err!=nil{
		MyLog.WarnLog.Println("模板渲染失败",err)
		return
	}
}
//接收上传的文件，并保存到数据库和本地
func acceptupdate(w http.ResponseWriter,r*http.Request){
	//获取id
	r.ParseForm()
	num,err:=strconv.Atoi(r.FormValue("workid"))
	if err!=nil{
		MyLog.WarnLog.Println("转换出错",err)
		return
	}
	//获取文件和内容
	r.ParseMultipartForm(32 << 20)
	name:=r.FormValue("name")
	//获取Work实例
	a,err:=sql.GetWork(num)
	if err!=nil{
		SetError(w,"内部错误")
		return
	}
	fileHeader:=r.MultipartForm.File["file"][0]
	file,err:=fileHeader.Open()
	date,err:=ioutil.ReadAll(file)
	if err!=nil{
		MyLog.WarnLog.Println("读取文件失败")
		SetError(w,"内部错误")
		return
	}
	err=sql.SetStudentWork(a.WorkID,name,date)
	t,err:=template.ParseFiles("./MyHandler/sucessUpdate.html")
	if err!=nil{
		MyLog.WarnLog.Println("打开模板出错")
		SetError(w,"内部错误")
		return
	}
	err=t.Execute(w,"")
	if err!=nil{
		MyLog.WarnLog.Println("模板渲染出错")
		SetError(w,"内部错误")
		return
	}
	MyLog.InfoLog.Println("收到上传的文件,文件名",name)
}