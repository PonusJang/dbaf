package models

import (
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"time"
)

type DbForward struct {
	gorm.Model
	Id         uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name       string    `json:"name" gorm:"dbName"`
	Type       int       `json:"dbType" gorm:"dbType"`
	ListenPort int       `json:"listenPort" gorm:"listenPort"`
	DbIp       string    `json:"dbIp" gorm:"dbIp"`
	DbPort     int       `json:"dbPort" gorm:"dbPort"`
	Created    time.Time `json:"created" gorm:"created"`
}
