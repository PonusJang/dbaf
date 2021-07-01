package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

type User struct {
	gorm.Model
	Id            int       `json:"id" gorm:"column:id;primaryKey"`
	Username      string    `json:"username" gorm:"column:username"`
	Password      string    `json:"password" gorm:"column:password"`
	Salt          string    `json:"salt" gorm:"column:salt"`
	Role          int       `json:"role" gorm:"column:roleId"`
	LastLogonIp   string    `json:"lastLogonIp" gorm:"column:lastLogonIp"`
	LastLogonDate time.Time `json:"lastLogonDate" gorm:"column:lastLogonDate"`
	CreatedAt     time.Time `json:"created" gorm:"column:created"`
}
