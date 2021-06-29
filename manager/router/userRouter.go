package router

import (
	"dbaf/manager/controllers"
	"github.com/gin-gonic/gin"
)

func LoadUserRouter(g *gin.Engine) {
	userRouter := g.Group("user")
	{
		//user.GET("/user/getList",)
		userRouter.POST("/user/create", controllers.CreateUser)
	}
}
