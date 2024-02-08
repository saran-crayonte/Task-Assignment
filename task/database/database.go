package database

import (
	"log"

	"github.com/saran-crayonte/task/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	conn := "host=localhost user=postgres password=1234 dbname=sample port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(conn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
		return
	}
	DB = db
	DB.AutoMigrate(&models.Task{})
	DB.AutoMigrate(&models.User{})
	DB.AutoMigrate(&models.Holiday{})
	DB.AutoMigrate(&models.TaskAssignment{})
}
