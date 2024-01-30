package database

import (
	"expensesManage/config"
	"expensesManage/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectWithDatabase() {
	database, err := gorm.Open(mysql.Open(config.BD_PATH), &gorm.Config{})
	if err != nil {
		panic("Could not connect to database")
	}
	DB = database
	DB.AutoMigrate(
		&models.UserBase{},
		&models.UserProfile{},
		&models.Categories{},
		&models.IncomeTable{},
		&models.Currency{},
		&models.ExpensesTable{})
}
