package controllers

import (
	"errors"
	"fmt"
	"isp/models"
	"log"
	"os"
	"unicode"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// SanitizeSignupData Sanitize the data passed in the signup form
func SanitizeSignupData(tempuser TempUser) (TempUser, Errors) {
	var errors Errors
	var err Err

	if tempuser.Gender != " " {
		err.message = "Gender is missing"
		errors.list = append(errors.list, err)
		return tempuser, errors
	}

	return tempuser, errors
}

// IsPasswordStrong Check password strength
func IsPasswordStrong(password string) (bool, error) {
	var IsLength, IsUpper, IsLower, IsNumber, IsSpecial bool

	if len(password) < 6 {
		return false, errors.New("Password Length should be more then 6")
	}
	IsLength = true

	for _, v := range password {
		switch {
		case unicode.IsNumber(v):
			IsNumber = true

		case unicode.IsUpper(v):
			IsUpper = true

		case unicode.IsLower(v):
			IsLower = true

		case unicode.IsPunct(v) || unicode.IsSymbol(v):
			IsSpecial = true

		}
	}

	if IsLength && IsLower && IsUpper && IsNumber && IsSpecial {
		return true, nil
	}

	return false, errors.New("Password validation failed.")

}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 8)
	if err != nil {
		log.Fatal("Error in Hashing")
		return "", err
	}
	return string(hashedPassword), err
}

func comparePasswords(hashedPwd string, plainPwd []byte) bool {
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}

func CheckCredentials(useremail, userpassword string, db *gorm.DB) bool {
	// db := c.MustGet("db").(*gorm.DB)
	// var db *gorm.DB
	var User models.User
	// Store user supplied password in mem map
	var expectedpassword string
	// check if the email exists
	err := db.Where("email = ?", useremail).First(&User).Error
	if err == nil {
		// User Exists...Now compare his password with our password
		expectedpassword = User.Password
		if err = bcrypt.CompareHashAndPassword([]byte(expectedpassword), []byte(userpassword)); err != nil {
			// If the two passwords don't match, return a 401 status
			log.Println("User is Not Authorized")
			return false
		}
		// User is AUthenticates, Now set the JWT Token
		fmt.Println("User Verified")
		return true
	} else {
		// returns an empty array, so simply pass as not found, 403 unauth
		log.Fatal("idhar ", err)

	}
	return false
}

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

	switch {
	case DB_HOST == "":
		fmt.Println("Environmet variable DB_HOST not set.")
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
	}

	return settings
}
