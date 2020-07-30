package dao

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"log"
)

var DB *gorm.DB

// Init 初始化mysql数据库
func Init() *gorm.DB {
	open, err := gorm.Open("mysql", "root:syouZX@tcp(localhost)/guguyun?charset=utf8&parseTime=true")
	if err != nil {
		log.Printf("sql init failured:%v", err)
		return DB
	}
	DB = open
	return DB
}
