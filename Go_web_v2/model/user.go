package model

import "gorm.io/gorm"

// 如果要用form-data发送请求，需要在user结构体定义时加上form的标签，这样才能绑定字段，不然会一直为空，
// 或者可以用raw的json发送请求。
// 不过我觉得shouldbind比较好
type User struct {
	gorm.Model
	Name      string `gorm:"varchar(20);not null"`
	Telephone string `gorm:"varchar(20);not null;unique"`
	// 未使用驼峰命名
	PassWord string `gorm:"varchar(20);not null"`
}
