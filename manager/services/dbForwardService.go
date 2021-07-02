package services

import (
	log "dbaf/log"
	db "dbaf/manager/databases"
	"dbaf/manager/models"
)

func AddDbForward(d *models.DbForward) bool {
	err := db.Db.Model(&models.DbForward{}).Create(d).Error
	if err != nil {
		log.Error(err)
		return false
	}
	return true
}

func DeleteForward(id int) bool {
	err := db.Db.Model(&models.DbForward{}).Delete(id).Error
	if err != nil {
		log.Error(err)
		return false
	}
	return true
}

func UpdateForward(id int, tmpDbForward *models.DbForward) bool {
	var origin models.DbForward
	var count int
	db.Db.Model(&models.DbForward{}).Find(&origin, id).Count(&count)
	if count == 0 {
		return false
	} else {
		err := db.Db.Model(&origin).Updates(&tmpDbForward).Error
		if err != nil {
			return false
		}
		return true
	}
}

func FindForwardByServer(server string) []models.DbForward {
	var tmp []models.DbForward
	db.Db.Model(&models.DbForward{}).Where("server = ?", server).Find(&tmp)
	return tmp
}

func GetDbForwardAll() (int, []models.DbForward) {
	var tmp []models.DbForward
	var count int
	db.Db.Model(&models.DbForward{}).Find(tmp).Count(&count)
	return count, tmp
}

func GetDbForwardList(page int, limit int, param map[string]interface{}) (int, []models.DbForward) {
	var tmp []models.DbForward
	var count int

	if len(param) > 0 {
		for k, v := range param {
			switch k {
			case "dbIp":
				if v != "" {
					db.Db.Where("dbIp = ?", v).Limit(limit).Offset(page * limit).Order("created desc").Find(&tmp).Count(&count)
					break
				}
			case "dbPort":
				if v != 0 {
					db.Db.Where("dbPort = ?", v).Limit(limit).Offset(page * limit).Order("created desc").Find(&tmp).Count(&count)
					break
				}
			case "listenPort":
				if v != 0 {
					db.Db.Where("listenPort = ?", v).Limit(limit).Offset(page * limit).Order("created desc").Find(&tmp).Count(&count)
					break
				}
			default:
				db.Db.Limit(limit).Offset(page * limit).Order("created desc").Find(&tmp).Count(&count)
			}
		}
	} else {
		db.Db.Limit(limit).Offset(page * limit).Order("created desc").Find(&tmp).Count(&count)
	}
	return count, tmp
}
