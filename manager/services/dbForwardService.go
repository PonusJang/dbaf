package services

import (
	log "dbaf/log"
	db "dbaf/manager/databases"
	"dbaf/manager/models"
	uuid "github.com/satori/go.uuid"
)

func AddDbForward(d *models.DbForward) bool {
	err := db.Db.Create(d).Error
	if err != nil {
		log.Error(err)
		return false
	}
	return true
}

func DeleteForward(id uuid.UUID) bool {
	err := db.Db.Delete(id).Error
	if err != nil {
		log.Error(err)
		return false
	}
	return true
}

func UpdateForward(id uuid.UUID) bool {
	err := db.Db.Update(id).Error
	if err != nil {
		log.Error(err)
		return false
	}
	return true
}

func FindForwardByServer(server string) bool {
	err := db.Db.Find(server).Error
	if err != nil {
		log.Error(err)
		return false
	}
	return true
}
