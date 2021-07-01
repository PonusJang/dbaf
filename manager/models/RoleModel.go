package models

import (
	"github.com/jinzhu/gorm"
)

type Role struct {
	gorm.Model
	Id   int    `json:"id" gorm:"column:id;primaryKey"`
	Name string `json:"name" gorm:"column:name"`
}
