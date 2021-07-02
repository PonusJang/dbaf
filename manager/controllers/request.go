package controllers

type UserReqForm struct {
	Username string `json:"username"  binding:"required"`
	Password string `json:"password"  binding:"required"`
}

type DbForwardForm struct {
	Name       string `json:"name"  binding:"required"`
	ListenPort int    `json:"listenPort"  binding:"required"`
	DbIp       string `json:"dbIp"  binding:"required"`
	DbPort     int    `json:"dbPort"  binding:"required"`
	DbType     int    `json:"dbType"  binding:"required"`
}
