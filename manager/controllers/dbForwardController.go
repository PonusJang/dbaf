package controllers

import (
	"dbaf/manager/models"
	"dbaf/manager/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

//添加

func AddDbForward(c *gin.Context) {
	name := c.PostForm("name")
	listenPort := c.PostForm("listenPort")
	dbIp := c.PostForm("dbIp")
	dbPort := c.PostForm("dbPort")
	dbType := c.PostForm("type")

	lport, _ := strconv.Atoi(listenPort)
	dport, _ := strconv.Atoi(dbPort)
	t, _ := strconv.Atoi(dbType)

	tmp := &models.DbForward{Name: name, ListenPort: lport, DbIp: dbIp, DbPort: dport, Type: t}
	status := services.AddDbForward(tmp)
	if status {
		c.JSON(http.StatusOK, Ret{CODE_SUCCESS, status, "", nil})
	}
	c.JSON(http.StatusOK, Ret{CODE_FALURE, status, "", nil})
}

//更新

func UpdateDbForward(c *gin.Context) {

	c.JSON(200, Ret{CODE_FALURE, false, "", nil})
}

//删除

func DeleteDbForward(c *gin.Context) {
	id := c.Query("id")
	iD, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusOK, Ret{CODE_PARAM_ERROR, true, MSG_PARAM_ERROR, nil})
	} else {
		success := services.DeleteForward(iD)
		if success {
			c.JSON(http.StatusOK, Ret{CODE_SUCCESS, true, MSG_SUCCESS, nil})
		} else {
			c.JSON(http.StatusOK, Ret{CODE_FALURE, false, MSG_FAILURE, nil})
		}
	}
}

//查找

func FindDbForwardByServer(c *gin.Context) {
	server := c.Query("server")
	data := services.FindForwardByServer(server)
	c.JSON(http.StatusOK, Ret{CODE_FALURE, false, "", data})
}

//获取列表

func GetDbForwardList(c *gin.Context) {

	page := c.DefaultQuery("page", "0")
	limit := c.DefaultQuery("page", "10")

	pageNo, err1 := strconv.Atoi(page)
	if err1 != nil {
		c.JSON(http.StatusBadRequest, Ret{CODE_FALURE, false, MSG_FAILURE, nil})
	}
	pageSize, err2 := strconv.Atoi(limit)
	if err2 != nil {
		c.JSON(http.StatusBadRequest, Ret{CODE_FALURE, false, MSG_FAILURE, nil})
	}

	var param map[string]interface{}
	listen := c.Query("listen")
	server := c.Query("server")
	port := c.Query("port")
	param["dbIp"] = server
	param["dbPort"] = port
	param["listenPort"] = listen
	count, data := services.GetDbForwardList(pageNo, pageSize, param)
	resData := ResData{Count: count, Data: data}
	c.JSON(http.StatusOK, Ret{CODE_FALURE, false, "", resData})
}
