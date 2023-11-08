package database

import (
	"camping-backend/models"
	"log"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DB struct {
	Conn *gorm.DB
}

var Database DB

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
	err = db.AutoMigrate(
		&models.User{},
		new(models.Camping),
		new(models.Amenity),
		new(models.Tag),
	)
	if err != nil {
		log.Fatal(err)
	}

	Database.Conn = db

}

func (d DB) FindByAmenityId(amenityId int) (*models.Amenity, error) {
	var amenity models.Amenity
	err := Database.Conn.First(&amenity, "id = ?", amenityId).Error
	if err != nil {
		return nil, err
	}
	return &amenity, nil
}

func (d DB) FindByCampingId(campingId int) (*models.Camping, error) {
	var camping models.Camping
	err := Database.Conn.First(&camping, "id = ?", campingId).Error
	if err != nil {
		return nil, err
	}
	return &camping, nil
}

func (d DB) FindByTagId(tagId int) (*models.Tag, error) {
	var tag models.Tag
	err := Database.Conn.First(&tag, "id = ?", tagId).Error
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

func (d DB) FindByUserId(userId int) (*models.User, error) {
	var user models.User
	err := Database.Conn.First(&user, "id = ?", userId).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
