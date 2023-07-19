package routes

import (
	"Practice/Go-web/Go_web_v2/controller"

	"github.com/gin-gonic/gin"
)

func CollectRoutes(engine *gin.Engine) *gin.Engine {
	engine.POST("/register", controller.Register)
	engine.POST("/login", controller.Logic)

	return engine
}
