package config

import (
	"testing"
)

func TestInit( t *testing.T){
	a:=config{
		DataBase: sqlmessage{"mysql", "root", "123456", "127.0.0.1:3306", "homework", "utf8"	},
		Teacher:  teacher{1, "UserWang", "pw123", "wang"},
		Server:   serverMassage{"127.0.0.1:8080"},
	}
	if a!=Config{
		t.Fatal("got ",Config," want ",a)
	}
}
