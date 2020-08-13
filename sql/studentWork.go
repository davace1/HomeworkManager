/*
作者：CZ
时间：2020.8.5
内容：studentWork的功能实现
*/
package sql

import (
	"errors"
	"homeworkManager/MyLog"
	"os"
	"strconv"
)
//上传的实例化变量，记录文件的id和文件路劲
type StudentWork struct {
	workID   int    //所提交的作业id
	fileName string //文件的路径 格式"xxx.zip"
	filePath string //文件的路径 注意转义字符影响\ 可食用/代替

}

//用于创建StudentWork 作业号，文件名称 文件内容
func SetStudentWork(workID int, fileName string, data []byte) error {
	path := "./homeworkFile/" + strconv.Itoa(workID) + "/" + fileName
	if judgeWork(workID) == false {
		//失败 返回error
		return errors.New("Can Find work [by cz] ")

	}

	//保存文件
	err := save(data, path)
	if err != nil {
		return err
	}
	//向数据库插入文件
	_, err = db.Exec("INSERT INTO studentwork value (NULL,?,?);",
		workID, path)
	if err != nil {
		MyLog.WarnLog.Println("向studentwork表插入数据错误:", err)
		return err
	}
	return nil
}

//将上传内容写入文件
func save(data []byte, path string) error {
	//创建文件夹
	file, err := os.Create(path)
	if err != nil {
		path = "." + path
		//用于测试
		file, err = os.Create(path)
		if err != nil {
			MyLog.WarnLog.Println("创建文件失败", err)
			return err
		}
	}
	defer file.Close()
	//写入数据
	_, err = file.Write(data)
	if err != nil {
		MyLog.WarnLog.Println("写入文件失败", err)
	}
	return err
}

//删除文件
func (s StudentWork) DeleteWork(fileID int) error {
	str := s.filePath
	//删除文件
	err := os.Remove(str)
	if err != nil {
		str = "." + str
		err = os.Remove(str)
		if err != nil {
			MyLog.WarnLog.Println("删除文件失败", err)
			return err
		}
	}
	//从数据库中删除
	_, err = db.Exec("DELETE FROM studentwork where fileid=? ;", fileID)
	if err != nil {
		MyLog.WarnLog.Println("数据库删除记录失败", err)
		return err
	}
	return nil
}

/*
//测试\与转义字符是否影响结果
//结论 有影响 建议使用'/'而非'\\'
func (s StudentWork)ShowPath(FileID int)string{
	ans,err:=db.Query("SELECT filePath FROM studentwork WHERE fileid=?;",FileID)
	if err!=nil{
		panic(err)
	}
	ans.Next()
	var str string
	ans.Scan(&str)
	return str
}
*/
