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
	Name      string `gorm:"vaarchar(20);not null"`
	Telephone string `gorm:"vaarchar(20);not null;unique"`
	PassWord  string `gorm:"size:225;not null"`
}

func main() {
	// 连接数据库
	db, err := gorm.Open(mysql.Open("root:@(127.0.0.1:3306)/web?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{NamingStrategy: schema.NamingStrategy{
		SingularTable: true, // 禁用表名复数形式，例如User的表名默认是users,
	}})

	// 判定错误
	if err != nil {
		panic("failed to connect database,err:" + err.Error())
	}

	// 自动迁移表
	db.AutoMigrate(&User{})

	// 结束关闭数据库
	sqldb, _ := db.DB()
	defer sqldb.Close()

	r := gin.Default()

	// 注册
	r.POST("/register", func(c *gin.Context) {

		// 接受参数
		name := c.PostForm("name")
		telephone := c.PostForm("telephone")
		password := c.PostForm("password")

		// 判定名字
		if len(name) == 0 {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"code":    http.StatusUnprocessableEntity,
				"message": "用户名不能为空",
			})
			return
		}

		// 判定电话
		if len(telephone) != 11 {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"code":    http.StatusUnprocessableEntity,
				"message": "手机号必需为11位",
			})
			return
		}

		// 判定密码
		if len(password) < 6 {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"code":    http.StatusUnprocessableEntity,
				"message": "密码不能少于6位",
			})
			return
		}

		// 判定电话不能重复使用
		var user User
		db.Where("telephone = ?", telephone).First(&user)
		if user.ID != 0 {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"code":    http.StatusUnprocessableEntity,
				"message": "用户已存在",
			})
			return
		}

		// 避免重名
		// db.Where("name = ?", name).Find(&user)
		// if user.ID != 0 {
		// 	c.JSON(http.StatusUnprocessableEntity, gin.H{
		// 		"code":    http.StatusUnprocessableEntity,
		// 		"message": "用户名已存在",
		// 	})
		// 	return
		// }

		//加密密码存储
		// 接受两个参数：一个字节切片，表示要加密的密码，以及一个整数，表示加密成本
		hasedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		// 一个字节切片，表示生成的哈希密码，以及一个错误值
		// 生成的哈希密码可以用于安全地存储用户密码

		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"code":    500,
				"message": "加密失败",
			})
			return
		}
		// 创建用户数据
		newUser := User{
			Name:      name,
			Telephone: telephone,
			PassWord:  string(hasedPassword),
		}
		db.Create(&newUser)

		c.JSON(http.StatusOK, gin.H{
			"code":   200,
			"mesage": "创建成功",
		})

	})

	// 登陆
	r.POST("/login", func(c *gin.Context) {
		telephone := c.PostForm("telephone")
		password := c.PostForm("password")

		if len(telephone) != 11 {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"code":    http.StatusUnprocessableEntity,
				"message": "手机号必须为11位",
			})
			return
		}

		if len(password) < 6 {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"code":    http.StatusUnprocessableEntity,
				"message": "密码不能少于6位",
			})
			return
		}

		// 判断手机是否存在
		var user User
		db.Where("telephone = ?", telephone).First(&user)
		if user.ID == 0 {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"code":    http.StatusUnprocessableEntity,
				"message": "用户不存在",
			})
			return
		}
		// 将存储的哈希密码与用户输入的密码进行比较，以验证用户身份。
		if err := bcrypt.CompareHashAndPassword([]byte(user.PassWord), []byte(password)); err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"code":    422,
				"message": "密码错误",
			})

		}
		c.JSON(http.StatusOK, gin.H{
			"code":   200,
			"mesage": "登陆成功",
		})

	})

	panic(r.Run(":9090"))
}

func InitDB() *gorm.DB {
	// driverName := "mysql"
	// host := "127.0.0.1"
	// post := "3306"
	// database := "web"
	// username := "root"
	// password := ""
	// charset := "utf8mb4"
	// args := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=%s&parseTime=True",
	// 	username,
	// 	password,
	// 	host,
	// 	post,
	// 	database,
	// 	charset)
	// db, err := gorm.Open(driverName, args)

	db, err := gorm.Open(mysql.Open("root:root1234@(127.0.0.1:13306)/db1?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{NamingStrategy: schema.NamingStrategy{
		SingularTable: true, // 禁用表名复数形式，例如User的表名默认是users,
	}})

	if err != nil {
		panic("failed to connect database,err:" + err.Error())
	}

	db.AutoMigrate(&User{})

	return db
}
