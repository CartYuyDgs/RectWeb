package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Userinfos struct {
	ID     uint
	Name   string
	Age    string
	Hobby  string
	Gender string
}

func main() {
	db := ConnectMysql()
	u1 := Userinfos{2, "aaa", "18", "pingpang", "aaa"}

	result := db.Create(&u1)
	fmt.Println(result.Error, result.DB())
	CloseMysql(db)
}

func ConnectMysql() *gorm.DB {
	db, err := gorm.Open("mysql", "root:199411@(127.0.0.1:3306)/user_infos?charset=utf8mb4&parseTime=true&loc=Local")
	if err != nil {
		panic(err)
	}

	//defer db.Close()

	db.AutoMigrate(&Userinfos{})

	//创建数据行

	return db
}

func CloseMysql(db *gorm.DB) {
	db.Close()
}
