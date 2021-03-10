package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	//"time"
)

type Userinfos struct {
	ID     uint
	Name   string
	Age    string
	Hobby  string
	Gender string
	time string
}

type Hostinfo struct {
	ID 	int64 `gorm:"primary_key" json:"id"`
	HostName string
	HostId	 string
	Nices []Nics	`gorm:"ForeignKey:host_id;AssociationForeignKey:id"`
}

type Nics struct {
	ID string
	NicName string
	Hostid  string	`json:"host_id"`
}

func main() {
	db := ConnectMysql()

	n1 := Nics{
		ID:      "01",
		NicName: "nic1",
		Hostid:  "1111101",
	}

	n2 := Nics{
		ID:      "02",
		NicName: "nic2",
		Hostid:  "1111102",
	}

	h1 := Hostinfo{
		ID:       0,
		HostName: "abc1",
		HostId:   "1111101",
		Nices:    []Nics{n1,n2},
	}


	res2 := db.Create(&n1)

	fmt.Println(res2.Error)
	res2 = db.Create(&n2)
	fmt.Println(res2.Error)
	res := db.Create(&h1)
	fmt.Println(res.Error)
	//u6 := Userinfos{6, "aaa6", "18", "pingpang", "aaa",time.Now().String()}
	//u7 := Userinfos{7, "aaa7", "18", "pingpang", "aaa",time.Now().String()}
	//u8 := Userinfos{8, "aaa8", "18", "pingpang", "aaa",time.Now().String()}
	//u9 := Userinfos{9, "aaa9", "18", "pingpang", "aaa",time.Now().String()}
	//
	////us := []Userinfos{u5, u2, u3, u4}
	//db.Create(&u6)
	//db.Create(&u7)
	//db.Create(&u8)
	//db.Create(&u9)

	//fmt.Println(time.Now().String())
	//SearchInfo(6,db)
	//SearchInfo(7,db)
	//SearchInfo(8,db)
	//SearchInfo(9,db)
	//
	//ud3 := Userinfos{
	//	ID: 6,
	//}
	//DeleteInfo(&ud3,db)

	CloseMysql(db)
}

func DeleteInfo(user *Userinfos,db *gorm.DB) error {
	result := db.Delete(user)
	if result.Error != nil {
		fmt.Println(result.Error)
	}
	return nil
}

func SearchInfo(id int, db *gorm.DB) {
	u2 := Userinfos{}
	db.Find(&u2, id)
	//fmt.Println(res)
	fmt.Println(u2)
	fmt.Println(u2.time)
}

func ConnectMysql() *gorm.DB {
	db, err := gorm.Open("mysql", "root:199411@(127.0.0.1:3306)/user_infos?charset=utf8mb4&parseTime=true&loc=Local")
	if err != nil {
		panic(err)
	}

	//defer db.Close()

	db.AutoMigrate(Hostinfo{})
	db.AutoMigrate(Nics{})

	//创建数据行

	return db
}

func CloseMysql(db *gorm.DB) {
	db.Close()
}
