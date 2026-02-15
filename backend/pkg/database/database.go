package database

import (
	"fmt"
	"log"
	"lifeline/internal/models"
	
	"gorm.io/driver/sqlite"   
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect() {

	db, err := gorm.Open(sqlite.Open("lifeline.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), 
	})

	if err != nil {
		log.Fatal("Failed to connect to database. \n", err)
	}

	fmt.Println(" SQLite Database Connected Successfully!")
	
	db.Exec("PRAGMA foreign_keys = ON")

	log.Println("Running Migrations...")
	err = db.AutoMigrate(
		&models.User{},
		&models.Role{},
		&models.Bank{},
		&models.Address{},
		&models.Photo{},
		&models.BloodGroup{},
		&models.Moderator{},
		&models.Sponsor{},
		// Moroccan Geo Data
		&models.Region{},
		&models.City{},
		&models.BloodRequest{},
		&models.Notification{},
	)

	if err != nil {
		log.Fatal("Migration Failed: ", err)
	}
	
	fmt.Println(" Database Migrated Successfully!")

	DB = db
}