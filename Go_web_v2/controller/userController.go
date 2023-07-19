package controller

import (
	"Practice/Go-web/Go_web_v2/common"
	"Practice/Go-web/Go_web_v2/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *gin.Context) {
	// 定义db
	db := common.GetDB()

	var requsetUser model.User
	c.Bind(&requsetUser)
	name := requsetUser.Name
	telephone := requsetUser.Telephone
	password := requsetUser.PassWord

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

	var user model.User
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

	newUser := model.User{
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

func Logic(c *gin.Context) {
	db := common.GetDB()

	var requsetUser model.User
	c.Bind(&requsetUser)
	telephone := requsetUser.Telephone
	password := requsetUser.PassWord

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

	var user model.User
	db.Where("telephone = ?", telephone).First(&user)
	if user.ID == 0 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"code":    422,
			"message": "该电话未注册或输入错误",
		})
		return
	}

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
}
