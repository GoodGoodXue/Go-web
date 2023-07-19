package routes

import (
	"Practice/Go-web/Go_web_m/database"
	"Practice/Go-web/Go_web_m/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// var ctx *gin.Context
var db = database.GetDB

func Register(c *gin.Context) {
	// 对应模型
	// 接受参数
	User := models.User
	name := models.User.Name
	telephone := models.User.Telephone
	password := models.User.PassWord

	user, err := c.ShouldBind(&models.User)
	if err != nil {
		return nil, err
	}

	// telephone := c.Bind(&models.User.Telephone)
	// password := c.Bind(&models.User.PassWord)

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

	var user models.User
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
}

func Login(c *gin.Context) {

}

func ApiRoutes(engin *gin.Engine) {
	engin.PSOT("/register", Register)
	engin.PSOT("/login", Login)
}
