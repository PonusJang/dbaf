package databases

import (
	"dbaf/manager/models"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var Db *gorm.DB

func init() {
	var err error
	Db, err = gorm.Open("mysql", "root:root@tcp(127.0.0.1:3306)/dbaf?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		fmt.Printf("mysql connect error %v", err)
	}
	if Db.Error != nil {
		fmt.Printf("database error %v", Db.Error)
	}

	Db.AutoMigrate(&models.User{})

	Db.AutoMigrate(&models.DbForward{})

	Db.AutoMigrate(&models.Role{})
}
