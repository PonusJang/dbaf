package models

import (
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"time"
)

type User struct {
	gorm.Model
	Id            uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Username      string    `json:"username" gorm:"column:username"`
	Password      string    `json:"password" gorm:"column:password"`
	Salt          string    `json:"salt" gorm:"column:salt"`
	Role          uuid.UUID `json:"role" gorm:"column:role"`
	LastLogonIp   string    `json:"lastLogonIp" gorm:"column:lastLogonIp"`
	LastLogonDate time.Time `json:"lastLogonDate" gorm:"column:lastLogonDate"`
	CreatedAt     time.Time `json:"created_at"`
}
