package router

import (
	"dbaf/manager/controllers"
	"github.com/gin-gonic/gin"
)

func LoadDbForwardRouter(g *gin.Engine) {
	userRouter := g.Group("api/v1/dbForward")
	{
		userRouter.POST("/add", controllers.AddDbForward)
		userRouter.POST("/update", controllers.UpdateDbForward)
		userRouter.GET("/delete", controllers.DeleteDbForward)
		userRouter.GET("/findByServer", controllers.FindDbForwardByServer)
		userRouter.GET("/getList", controllers.GetDbForwardList)
		userRouter.GET("/getAll", controllers.GetDbForwardAll)
	}
}
