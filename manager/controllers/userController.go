package controllers

import (
	"dbaf/manager/models"
	"dbaf/manager/services"
	"github.com/gin-gonic/gin"
)

func CreateUser(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	u := &models.User{Username: username, Password: password}
	if services.Create(u) {
		c.JSON(200, &Ret{CODE_SUCCESS, true, MSG_SUCCESS, ""})
	} else {
		c.JSON(200, &Ret{CODE_FALURE, false, MSG_FAILURE, ""})
	}
}

func Login(c *gin.Context) {

}
