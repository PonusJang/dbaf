package controllers

import (
	"dbaf/manager/models"
	"dbaf/manager/services"
	"github.com/gin-gonic/gin"
	"strconv"
)

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
		c.JSON(200, &Ret{CODE_SUCCESS, status, "", nil})
	}
	c.JSON(200, &Ret{CODE_FALURE, status, "", nil})
}

func UpdateDbForward(c *gin.Context) {

	c.JSON(200, &Ret{CODE_FALURE, false, "", nil})
}

func DeleteDbForward(c *gin.Context) {

	c.JSON(200, &Ret{CODE_FALURE, false, "", nil})
}

func FindDbForward(c *gin.Context) {

	c.JSON(200, &Ret{CODE_FALURE, false, "", nil})
}

func GetListDbForward(c *gin.Context) {

	c.JSON(200, &Ret{CODE_FALURE, false, "", nil})
}