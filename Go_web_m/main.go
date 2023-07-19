package main

import (
	"Practice/Go-web/Go_web_m/database"
	"Practice/Go-web/Go_web_m/models"
	"Practice/Go-web/Go_web_m/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	engin := gin.Default()

	database.Connect()

	db := database.GetDB()
	db.AutoMigrate(&models.User{})

	routes.ApiRoutes(engin)
}
