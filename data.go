package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB

func DatabaseConnect(){
	var err error
	username := "MainA"
	password := "Vera1225"
	endpoint := "tcp(mydatabase.cmnihacpohfv.ap-southeast-2.rds.amazonaws.com)"
	dbname := "Trial1"
	params := "charset=utf8&parseTime=True&loc=Local"

	dataSource := username + ":" + password + "@" + endpoint + "/" + dbname + "?" + params

	db, err = gorm.Open("mysql", dataSource)
	if err != nil {
		panic(err)
	}

	//// Migrate the schema
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Post{})

}

//// Update - update product's price to 2000
//db.Model(&User).Update("Price", 2000)
//
//// Delete - delete product
//db.Delete(&User)
