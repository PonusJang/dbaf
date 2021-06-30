package manager

import (
	_ "dbaf/manager/databases"
	"dbaf/manager/router"
	"github.com/gin-gonic/gin"
	"net/http"
)

// 管理 web端

func RunServer() {
	gin.SetMode(gin.DebugMode)
	r := gin.Default()
	router.LoadUserRouter(r)

	http.ListenAndServe(":8835", r)
}
