package controllers

import (
	"dbaf/manager/models"
	"dbaf/manager/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

//添加用户

func CreateUser(c *gin.Context) {

	var userReq UserReqForm
	err := c.BindJSON(&userReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, Ret{CODE_PARAM_ERROR, false, MSG_PARAM_ERROR, nil})
		return
	}
	if userReq.Username == "" || userReq.Password == "" {
		c.JSON(http.StatusLengthRequired, Ret{CODE_PARAM_ERROR, false, MSG_PARAM_ERROR, nil})
		return
	}
	pass := []byte(userReq.Password)
	u := &models.User{Username: userReq.Username, Password: pass, CreatedAt: time.Now()}
	if services.CreateUser(u) {
		c.JSON(http.StatusOK, Ret{CODE_SUCCESS, true, MSG_SUCCESS, nil})
		return
	}
	c.JSON(http.StatusOK, Ret{CODE_FALURE, false, MSG_FAILURE, ""})
}

//登录

func Login(c *gin.Context) {

	var userReq UserReqForm
	err := c.BindJSON(&userReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, Ret{CODE_PARAM_ERROR, false, MSG_PARAM_ERROR, nil})
		return
	}
	if userReq.Username == "" || userReq.Password == "" {
		c.JSON(http.StatusLengthRequired, Ret{CODE_PARAM_ERROR, false, MSG_PARAM_ERROR, nil})
		return
	}

	loginIp := c.Request.RemoteAddr
	pass := []byte(userReq.Password)
	u := &models.User{Username: userReq.Username, Password: pass, LastLogonIp: loginIp, LastLogonDate: time.Now()}
	token, err := services.Login(u)
	if err != nil || token == "" {
		c.JSON(http.StatusOK, Ret{CODE_LOGIN_FAILURE, false, MSG_FAILURE, nil})
	} else {
		data := Token{Token: token}
		c.JSON(http.StatusOK, Ret{CODE_LOGIN_SUCCESS, true, MSG_LOGIN_SUCCESS, data})
	}
}

//删除用户

func DeleteUser(c *gin.Context) {
	id := c.Query("id")
	iD, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusOK, Ret{CODE_FALURE, false, MSG_FAILURE, nil})
	}
	success := services.DeleteUser(iD)
	if success {
		c.JSON(http.StatusOK, Ret{CODE_FALURE, false, MSG_FAILURE, nil})
	} else {
		c.JSON(http.StatusOK, Ret{CODE_LOGIN_SUCCESS, true, MSG_SUCCESS, nil})
	}
}

func UpdateUser(c *gin.Context) {
	id := c.PostForm("id")
	iD, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusOK, Ret{CODE_FALURE, false, MSG_FAILURE, nil})
	}
	var param map[string]interface{}
	password := c.PostForm("password")
	param["password"] = password

	if services.UpdateUser(iD, param) {
		c.JSON(http.StatusOK, Ret{CODE_FALURE, false, MSG_FAILURE, nil})
	} else {
		c.JSON(http.StatusOK, Ret{CODE_LOGIN_SUCCESS, true, MSG_SUCCESS, nil})
	}
}

func GetUserList(c *gin.Context) {

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
	username := c.Query("username")
	if username != "" {
		param["username"] = username
	}
	count, data := services.GetUserList(pageNo, pageSize, param)
	resData := ResData{Count: count, Data: data}
	c.JSON(http.StatusOK, Ret{CODE_FALURE, false, "", resData})
}

func FindByUsername(c *gin.Context) {

}
