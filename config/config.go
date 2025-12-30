package config

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	DSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", Mysql.USERNAME, Mysql.PASSWORD, Mysql.HOST, Mysql.PORT, Mysql.NAME)
	DBtemp, err := gorm.Open(mysql.Open(DSN), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	DB = DBtemp
	fmt.Println("Connected to MySQL")
}
