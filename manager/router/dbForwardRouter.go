package router

import (
	"dbaf/manager/controllers"
	"github.com/gin-gonic/gin"
)

func LoadDbForwardRouter(g *gin.Engine) {
	userRouter := g.Group("dbForward")
	{
		userRouter.POST("/dbForward/add", controllers.AddDbForward)
		userRouter.POST("/dbForward/update", controllers.UpdateDbForward)
		userRouter.GET("/dbForward/delete", controllers.DeleteDbForward)
		userRouter.GET("/dbForward/findByServer", controllers.FindDbForwardByServer)
		userRouter.GET("/dbForward/getList", controllers.GetDbForwardList)
	}
}
