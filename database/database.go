package database

import (
	"camping-backend/models"
	"log"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDB() {
	var DSN string = "sqlite.db"
	db, err := gorm.Open(sqlite.Open(DSN), &gorm.Config{
		NowFunc: func() time.Time {
			return time.Now().Local()
		},
	})

	if err != nil {
		log.Fatal("Failed to connect to the database")
	}

	log.Println("Connected to the database successfully")
	db.Logger = logger.Default.LogMode(logger.Info)
	log.Println("Running Migrations")

	// Todo: Add migrations
	err = db.AutoMigrate(&models.User{}, new(models.Camping))
	if err != nil {
		log.Fatal(err)
	}

	DB = db

}
