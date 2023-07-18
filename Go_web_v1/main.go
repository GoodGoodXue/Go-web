package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type User struct {
	gorm.Model
	Name      string `gorm:"varchar(20);not null"`
	Telephone string `gorm:"varchar(20);not null;unique"`
	// 未使用驼峰命名
	PassWord string `gorm:"varchar(20);not null"`
}

func main() {

	db, err := gorm.Open(mysql.Open("root:@(127.0.0.1:3306)/web?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 禁用表名复数形式，例如User的表名默认是users,
		}})

	// 错误信息的内容取决于 err 变量的值、
	// 后半部分为 err.Error() 的返回值，即 err 变量所表示的错误信息
	if err != nil {
		panic("failed to connect database,err" + err.Error())
	}

	db.AutoMigrate(&User{})

	sqldb, _ := db.DB()
	defer sqldb.Close()

	r := gin.Default()

	r.POST("/register", func(c *gin.Context) {

		name := c.PostForm("name")
		telephone := c.PostForm("telephone")
		password := c.PostForm("password")

		if len(name) == 0 {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"code":    422,
				"message": "用户名不能为空",
			})
			return
		}

		if len(telephone) != 11 {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"code":    422,
				"message": "电话必须为11位",
			})
			return
		}

		if len(password) < 6 {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"code":    422,
				"message": "密码不能小于6位",
			})
			return
		}

		var user User
		db.Where("telephone = ?", telephone).First(&user)
		if user.ID != 0 {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"code":    422,
				"message": "电话已注册",
			})
			return
		}

		hasedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"code":    500,
				"message": "密码加密失败",
			})
		}

		newUser := User{
			Name:      name,
			Telephone: telephone,
			PassWord:  string(hasedPassword),
		}
		db.Create(&newUser)
		// fmt.Println(NewUser)

		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"message": "注册成功",
		})

	})

	r.POST("/login", func(c *gin.Context) {

		// name := c.PostForm("name")
		telephone := c.PostForm("telephone")
		password := c.PostForm("password")

		var user User
		// db.Where("name = ?", name).First(&user)
		// if user.ID == 0 {
		// 	c.JSON(http.StatusUnprocessableEntity, gin.H{
		// 		"code":    422,
		// 		"message": "未找到该用户名",
		// 	})
		// 	return
		// }

		// 错误思维
		// 统一判断完电话再判定其他，容易造成资源浪费

		// 优先简单判断传入格式，再判定其他
		if len(telephone) != 11 {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"code":    422,
				"message": "号码必须为11位",
			})
			return
		}

		if len(password) < 6 {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"code":    422,
				"message": "密码不能少于6位",
			})
			return
		}

		db.Where("telephone = ?", telephone).First(&user)
		if user.ID == 0 {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"code":    422,
				"message": "该电话未注册或输入错误",
			})
			return
		}

		// 错误思维
		// 通过用户输入的密码进行查找，匹配数据库中的加密密码，但是未经过加密无法匹配
		// err := db.Model(&User{}).Where("password = ?", password).First(&user)
		// 一个参数数据库加密密码，第二参数为传入的密码
		// err := bcrypt.CompareHashAndPassword(MiPassword, []byte(password))

		// if user.ID == 0 {
		// 	c.JSON(http.StatusUnprocessableEntity, gin.H{
		// 		"code":    422,
		// 		"message": "密码错误",
		// 	})
		// 	return
		// }

		// 用于比较哈希密码和明文密码，接受两个参数：一个字节切片，表示哈希密码，以及一个字节切片，表示明文密码
		if err := bcrypt.CompareHashAndPassword([]byte(user.PassWord), []byte(password)); err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"code":    422,
				"message": "密码错误",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "登陆成功",
		})

	})

	// 最后监听
	panic(r.Run(":9090"))

}

// func InitDB() *gorm.DB {
// 	driverName := "mysql"
// 	host := "127.0.0.1"
// 	post := "3306"
// 	database := "web"
// 	username := "root"
// 	password := ""
// 	charset := "utf8mb4"
// 	args := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=%s&parseTime=True",
// 		username,
// 		password,
// 		host,
// 		post,
// 		database,
// 		charset)

// 	db, err := gorm.Open(driverName, args)
// 	if err != nil {
// 		panic("failed to connect database,err:" + err.Error())
// 	}

// 	// db, err := gorm.Open(mysql.Open("root:root1234@(127.0.0.1:13306)/db1?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{NamingStrategy: schema.NamingStrategy{
// 	// 	SingularTable: true, // 禁用表名复数形式，例如User的表名默认是users,
// 	// }})

// 	db.AutoMigrate(&User{})

// 	return db
// }
