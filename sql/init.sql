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