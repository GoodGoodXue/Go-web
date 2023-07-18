package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/cors/wrapper/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type User struct {
	gorm.Model
	Name      string `gorm:"vaarchar(20);not null"`
	Telephone string `gorm:"vaarchar(20);not null;unique"`
	Password  string `gorm:"vaarchar(20);not null"`
}

func main() {
	r := gin.Default()

	db, err := gorm.Open(mysql.Open("root:@(127.0.0.1:3306)/web?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 禁用表名复数形式，例如User的表名默认是users,
		}})

	if err != nil {
		panic(err) // "failed to database:",
	}

	sqldb, _ := db.DB()
	defer sqldb.Close()

	r.Run(":9090")

	db.AutoMigrate(&User{})

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

		if len(password) < 5 {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"code":    422,
				"message": "密码必须大于6位",
			})
			return
		}

		MiPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"code":    422,
				"message": "密码加密失败",
			})
		}

		var user User
		user = User{
			Name:      name,
			Telephone: telephone,
			Password:  string(MiPassword),
		}
		db.Create(&user)

		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"message": "注册成功",
		})

	})

	r.POST("/Login", func(c *gin.Context) {

		name := c.PostForm("name")
		telephone := c.PostForm("telephone")
		password := c.PostForm("password")

		var user User
		db.Model(&User{}).Where("name = ?", name).First(&user)
		if user.ID == 0 {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"code":    422,
				"message": "未找到该用户名",
			})
			return
		}

		db.Model(&User{}).Where("telephone = ?", telephone).First(&user)
		if user.ID == 0 {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"code":    422,
				"message": "该电话未注册或输入错误",
			})
			return
		}

		// err := db.Model(&User{}).Where("password = ?", password).First(&user)
		err := bcrypt.CompareHashAndPassword(MiPassword, []byte(password))

		if user.ID == 0 {
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

}
