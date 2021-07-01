package router

import (
	"dbaf/manager/controllers"
	"github.com/gin-gonic/gin"
)

func LoadUserRouter(g *gin.Engine) {
	userRouter := g.Group("api/v1/user")
	{
		userRouter.POST("/login", controllers.Login)
		userRouter.POST("/create", controllers.CreateUser)
		userRouter.GET("/delete", controllers.DeleteUser)
		userRouter.POST("/update", controllers.UpdateUser)
		userRouter.GET("/getList", controllers.GetUserList)
		userRouter.GET("/findByUsername", controllers.FindByUsername)
	}
}
