# HomeworkManager-
使用go和MySql的简单的收作业网页

## 第三方库：

### [go-sql-driver](https://github.com/go-sql-driver)/：

用于连接MySQL

安装：

```go
$ go get -u github.com/go-sql-driver/mysql
```

**修改配置文件** `config.json`

####执行数据库初始化

使用MySQL，请根据需要修改

```mysql
CREATE DATABASE homework;
USE homework;
CREATE TABLE work(
    workID INT PRIMARY KEY auto_increment,
    teacherID INT,
    name VARCHAR(50) NOT NULL,
    endtime DATETIME NOT NULL,
    begintime DATETIME NOT NULL,
    INDEX name_index (name)
)engine=InnoDB charset='utf8';
CREATE TABLE studentWork (
    fileID  INT PRIMARY KEY auto_increment,
    workID INT NOT NULL,
    filePath VARCHAR(50) NOT NULL,
    CONSTRAINT workID_foreign FOREIGN KEY(workID) REFERENCES work(workID) ON DELETE CASCADE ON UPDATE CASCADE
)engine=InnoDB charset='utf8';
```

###文件存储
项目文件夹下的`logFile`用于存储日志文件

项目文件下的`homeworkFile`用于存储上传的文件

项目文件夹下的`ZIP`用于存储压缩的文件夹
