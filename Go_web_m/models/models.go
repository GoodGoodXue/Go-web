package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name      string `gorm:"varchar(20);not null"`
	Telephone string `gorm:"varchar(20);not null;unique"`
	// 未使用驼峰命名
	PassWord string `gorm:"varchar(20);not null"`
}
