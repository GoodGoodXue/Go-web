package database

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB

func Connect() (db *gorm.DB) {
	db, err := gorm.Open(mysql.Open("root:@(127.0.0.1:3306)/web?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 禁用表名复数形式，例如User的表名默认是users,
		}})

	// 错误信息的内容取决于 err 变量的值、
	// 后半部分为 err.Error() 的返回值，即 err 变量所表示的错误信息
	if err != nil {
		panic("failed to connect database,err" + err.Error())
	}

	sqldb, _ := db.DB()
	defer sqldb.Close()

	DB = db

	return

}

func GetDB() *gorm.DB {
	return DB
}
