package services

import (
	"dbaf/manager/common"
	db "dbaf/manager/databases"
	"dbaf/manager/models"
	uuid "github.com/satori/go.uuid"
)

func Login(u *models.User) (string, error) {
	var tmpUser *models.User
	db.Db.First(tmpUser, u.Username)
	if tmpUser.Password == u.Password {
		return common.GenerateToken(u.Username, u.Password)
	} else {
		return "", nil
	}
}

func Verify(token string) (string, error) {
	claims, err := common.ParseToken(token)
	if err != nil {
		return "", nil
	}
	return claims.Username, nil
}

func Create(u *models.User) bool {
	err := db.Db.Create(u).Error
	if err != nil {
		return false
	} else {
		return true
	}
}

func Delete(id uuid.UUID) bool {
	err := db.Db.Delete(&models.User{}, id).Error
	if err != nil {
		return false
	}
	return true
}

func Update(id uuid.UUID, u *models.User) bool {
	tmp := db.Db.First(&models.User{}, id)
	err := db.Db.Model(u).Updates(&tmp).Error
	if err != nil {
		return false
	}
	return true
}
