package controllers

import (
	logger "dbaf/log"
	"dbaf/manager/models"
	"dbaf/manager/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type Token struct {
	token string `json:"token"`
}

type UserReq struct {
	username string `json:"username"  binding:"required"`
	password string `json:"password"  binding:"required"`
}

//添加用户

func CreateUser(c *gin.Context) {
	var  userReq UserReq
	err := c.ShouldBindJSON(&userReq); if err!=nil {
		c.JSON(http.StatusBadRequest, Ret{CODE_SUCCESS, true, MSG_SUCCESS, nil})
		return
	}
	logger.Debug(c.Request.Body.)
	logger.Info(userReq)
	logger.Info(userReq.username)
	logger.Info(userReq.password)
	u := &models.User{Username: userReq.username, Password: userReq.password}
	if services.CreateUser(u) {
		c.JSON(http.StatusOK, Ret{CODE_SUCCESS, true, MSG_SUCCESS, nil})
		return
	}
	c.JSON(http.StatusOK, Ret{CODE_FALURE, false, MSG_FAILURE, ""})
}

//登录

func Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	loginIp := c.Request.RemoteAddr
	u := &models.User{Username: username, Password: password,LastLogonIp: loginIp,LastLogonDate: time.Now()}
	token, err := services.Login(u)
	if err != nil {
		c.JSON(http.StatusOK, Ret{CODE_LOGIN_FAILURE, false, MSG_FAILURE, nil})
	} else {
		data := Token{token: token}
		c.JSON(http.StatusOK, Ret{CODE_LOGIN_SUCCESS, true, MSG_SUCCESS, data})
	}
}

//删除用户

func DeleteUser(c *gin.Context) {
	id := c.Query("id")
	iD,err := strconv.Atoi(id)
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


func UpdateUser(c *gin.Context){
	id := c.PostForm("id")
	iD,err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusOK, Ret{CODE_FALURE, false, MSG_FAILURE, nil})
	}
	var param map[string]interface{}
	password := c.PostForm("password")
	param["password"] = password

	if services.UpdateUser(iD,param){
		c.JSON(http.StatusOK, Ret{CODE_FALURE, false, MSG_FAILURE, nil})
	}else {
		c.JSON(http.StatusOK, Ret{CODE_LOGIN_SUCCESS, true, MSG_SUCCESS, nil})
	}
}

func GetUserList(c *gin.Context)  {

	page := c.DefaultQuery("page","0")
	limit := c.DefaultQuery("page","10")

	pageNo,err1 := strconv.Atoi(page)
	if err1!=nil{
		c.JSON(http.StatusBadRequest, Ret{CODE_FALURE, false, MSG_FAILURE, nil})
	}
	pageSize,err2 :=strconv.Atoi(limit)
	if err2!=nil{
		c.JSON(http.StatusBadRequest, Ret{CODE_FALURE, false, MSG_FAILURE, nil})
	}

	var param map[string]interface{}
	username := c.Query("username")
	if username!="" {
		param["username"] = username
	}
	count,data := services.GetUserList(pageNo, pageSize,param)
	resData := ResData{Count: count,Data: data}
	c.JSON(http.StatusOK, Ret{CODE_FALURE, false, "", resData})
}

func FindByUsername(c *gin.Context)  {

}