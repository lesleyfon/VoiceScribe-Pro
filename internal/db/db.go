package db

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type NotesStruct struct {
	gorm.Model
	Title string `gorm:"type:text"`
}

func (NotesStruct) TableName() string {
	return "notes"
}

func InitDB() (*gorm.DB, error) {
	DATABASE_URL := os.Getenv("POSTGRES_URL")

	if DATABASE_URL == "" {
		log.Fatal("No Database URL")
		return nil, nil
	}

	// Configure GORM
	config := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error), // Only log errors
	}

	db, err := gorm.Open(postgres.Open(DATABASE_URL), config)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	// Set reasonable connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if os.Getenv("ENV") == "development" {
		if err := db.AutoMigrate(&NotesStruct{}); err != nil {
			log.Printf("Migration error: %v", err)
		}
	}
	// Run AutoMigrate

	return db, nil
}
