package controllers

import (
	"dbaf/manager/models"
	"dbaf/manager/services"
	"github.com/gin-gonic/gin"
)

type Token struct {
	token string `json:"token"`
}

func CreateUser(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	u := &models.User{Username: username, Password: password}
	if services.CreateUser(u) {
		c.JSON(200, Ret{CODE_SUCCESS, true, MSG_SUCCESS, ""})
	} else {
		c.JSON(200, Ret{CODE_FALURE, false, MSG_FAILURE, ""})
	}
}

func Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	u := &models.User{Username: username, Password: password}
	token, err := services.Login(u)
	if err != nil {
		c.JSON(200, Ret{CODE_LOGIN_FAILURE, false, MSG_FAILURE, nil})
	} else {
		data := Token{token: token}
		c.JSON(200, Ret{CODE_LOGIN_SUCCESS, true, MSG_SUCCESS, data})
	}
}

func DeleteUser(c *gin.Context) {
	username := c.Query("username")
	success := services.DeleteUser(username)
	if success {
		c.JSON(200, Ret{CODE_FALURE, false, MSG_FAILURE, nil})
	} else {
		c.JSON(200, Ret{CODE_LOGIN_SUCCESS, true, MSG_SUCCESS, nil})
	}
}
