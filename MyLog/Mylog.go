package MyLog

import (
	"log"
	"os"
	"time"
)

var ErrorLog *log.Logger
var WarnLog *log.Logger
var InfoLog *log.Logger

var file *os.File

func init() {
	filename := ".\\logFile\\Log" + time.Now().Format("20060102150405") + ".log"
	var err error
	file, err = os.Create(filename)
	if err != nil {
		filename = "." + filename
		file, err = os.Create(filename)
		if err != nil {
			panic(err)
		}
	}
	//defer file.Close()
	//debug:文件关闭之后就写入失败了，请不要关闭
	ErrorLog = log.New(file, "error :", log.LstdFlags|log.Llongfile)
	WarnLog = log.New(file, "warn :", log.LstdFlags|log.Llongfile)
	InfoLog = log.New(file, "info :", log.LstdFlags|log.Llongfile)
}

//用于关闭日志文件
func close() error {
	err := file.Close()
	if err != nil {
		panic(err)
	}
	return err
}
