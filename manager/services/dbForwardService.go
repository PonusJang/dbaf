package services

import (
	log "dbaf/log"
	db "dbaf/manager/databases"
	"dbaf/manager/models"
)

func AddDbForward(d *models.DbForward) bool {
	err := db.Db.Create(d).Error
	if err != nil {
		log.Error(err)
		return false
	}
	return true
}
