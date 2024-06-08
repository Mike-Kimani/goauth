package database

import (
	"github.com/mike-kimani/goauth/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	connection, err := gorm.Open(mysql.Open("root:novek@/jwt_go_auth"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	DB = connection
	err = connection.AutoMigrate(&models.User{})
	if err != nil {

	}
}
