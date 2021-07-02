package services

import (
	"bytes"
	logger "dbaf/log"
	"dbaf/manager/common"
	db "dbaf/manager/databases"
	"dbaf/manager/models"
	"github.com/tjfoc/gmsm/sm3"
	"time"
)

func Login(u *models.User) (string, error) {
	var tmpUser models.User
	db.Db.Model(&models.User{}).Where("username = ?", u.Username).First(&tmpUser)
	h := sm3.New()
	logger.Debug(tmpUser.Password)
	logger.Debug(u.Password)
	logger.Debug(tmpUser.Salt)
	h.Write(append(u.Password, tmpUser.Salt...))
	if bytes.Equal(tmpUser.Password, h.Sum(nil)) {
		db.Db.Model(&models.User{}).Where("username = ?", u.Username).Update("lastLogonIp", u.LastLogonIp).Update("lastLogonDate", u.LastLogonDate)
		return common.GenerateToken(u.Username, string(u.Password))
	} else {
		return "", nil
	}
}

func IsExpires(token string) bool {
	claims, err := common.ParseToken(token)
	if err != nil {
		return false
	}
	if claims.VerifyExpiresAt(time.Now().Unix(), true) {
		return true
	} else {
		return false
	}
}

func VerifyToken(token string) (string, error) {
	claims, err := common.ParseToken(token)
	if err != nil {
		return "", nil
	}
	return claims.Username, nil
}

func CreateUser(u *models.User) bool {
	var count int
	logger.Debug(u.Username)
	db.Db.Model(&models.User{}).Where("username = ?", u.Username).Count(&count)
	logger.Debug(count)
	if count > 0 {
		return false
	}
	u.Salt = []byte(GetRandomString(8))
	h := sm3.New()
	h.Write(append(u.Password, u.Salt...))
	u.Password = h.Sum(nil)
	err := db.Db.Create(u).Error
	if err != nil {
		return false
	} else {
		return true
	}
}

func DeleteUser(id int) bool {
	err := db.Db.Delete(&models.User{}, id).Error
	if err != nil {
		return false
	}
	return true
}

func UpdateUser(id int, param map[string]interface{}) bool {
	err := db.Db.Model(&models.User{}).Where(id).Updates(param).Error
	if err != nil {
		return false
	}
	return true
}

func GetUserList(page int, limit int, param map[string]interface{}) (int, []models.User) {
	var tmp []models.User
	var count int

	if len(param) > 0 {
		v := param["username"]
		username := v.(string)
		db.Db.Limit(limit).Offset(page*limit).Where("username like ?", "%"+username+"%").Find(&tmp).Count(&count)
	} else {
		db.Db.Limit(limit).Offset(page * limit).Order("created desc").Find(&tmp).Count(&count)
	}
	return count, tmp
}
