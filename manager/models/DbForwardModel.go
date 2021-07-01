package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

type DbForward struct {
	gorm.Model
	Id         int       `json:"id" gorm:"column:id;primaryKey"`
	Name       string    `json:"name" gorm:"column:dbName"`
	Type       int       `json:"dbType" gorm:"column:dbType"`
	ListenPort int       `json:"listenPort" gorm:"column:listenPort"`
	DbIp       string    `json:"dbIp" gorm:"column:dbIp"`
	DbPort     int       `json:"dbPort" gorm:"column:dbPort"`
	Created    time.Time `json:"created" gorm:"column:created"`
}
