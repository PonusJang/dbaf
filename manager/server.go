package manager

import (
	_ "dbaf/manager/databases"
	"dbaf/manager/middlewares"
	"dbaf/manager/router"
	"github.com/gin-gonic/gin"
	"net/http"
)

// 管理 web端

func RunServer() {
	gin.SetMode(gin.DebugMode)
	r := gin.Default()
	r.Use(middlewares.LogerMiddleware())
	router.LoadUserRouter(r)
	router.LoadDbForwardRouter(r)
	http.ListenAndServe(":8835", r)
}
