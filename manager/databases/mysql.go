package databases

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var Db *gorm.DB

func init() {
	var err error
	Db, err = gorm.Open("mysql", "root:Y1z@2021@tcp(10.10.8.188:3307)/dbaf?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		fmt.Printf("mysql connect error %v", err)
	}
	if Db.Error != nil {
		fmt.Printf("database error %v", Db.Error)
	}
}
