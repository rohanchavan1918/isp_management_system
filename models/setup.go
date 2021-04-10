package models

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

type Settings struct {
	DB_HOST     string
	DB_NAME     string
	DB_USER     string
	DB_PASSWORD string
	DB_PORT     string
}

func InitializeSettings() Settings {
	DB_HOST := os.Getenv("DB_HOST")
	DB_NAME := os.Getenv("DB_NAME")
	DB_USER := os.Getenv("DB_USER")
	DB_PASSWORD := os.Getenv("DB_PASSWORD")
	DB_PORT := os.Getenv("DB_PORT")

	switch {
	case DB_HOST == "":
		fmt.Println("1 Environmet variable DB_HOST not set.")
		os.Exit(1)
	case DB_NAME == "":
		fmt.Println("Environmet variable DB_NAME not set.")
		os.Exit(1)
	case DB_USER == "":
		fmt.Println("Environmet variable DB_USER not set.")
		os.Exit(1)
	case DB_PASSWORD == "":
		fmt.Println("Environmet variable DB_PASSWORD not set.")
		os.Exit(1)
	}

	settings := Settings{
		DB_HOST:     DB_HOST,
		DB_NAME:     DB_NAME,
		DB_USER:     DB_USER,
		DB_PASSWORD: DB_PASSWORD,
		DB_PORT:     DB_PORT,
	}

	return settings
}

func ConnectDataBase() {
	// THis file is used to initialize the defined models
	settings := InitializeSettings()
	// dsn := "host=localhost user=postgres password=postgres dbname=goauth port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	dsn := "host=" + settings.DB_HOST + " user=" + settings.DB_USER + " password=" + settings.DB_PASSWORD + " dbname=" + settings.DB_NAME + " port=" + settings.DB_PORT + " sslmode=disable TimeZone=Asia/Kolkata"
	fmt.Println(dsn)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("Failed to connect to database!")
	}

	db.AutoMigrate(&User{}, &ForgotPassword{})

	DB = db
}
