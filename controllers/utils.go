package controllers

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"html/template"
	"isp/models"
	"log"
	"math/big"
	"net/smtp"
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
		log.Fatal("ERR ", err)

	}
	return false
}

// GenerateSecureToken returns a secured n digit token
func GenerateSecureToken(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}

func getRandNum() int64 {
	max := big.NewInt(999999)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		log.Fatal(err)
	}
	otp := n.Int64()
	return otp
}

func GetOtp() int64 {
	for {
		otp := getRandNum()

		if otp > 100000 && otp < 999999 {
			return otp
		}
	}
}

// SendMail - GLobal function to send mail
func SendMailToUser(recipentList []string, otp string, purpose string) {
	SMTP_USER := os.Getenv("SMTP_USER")
	SMTP_PASSWORD := os.Getenv("SMTP_PASSWORD")
	SMTP_HOST := os.Getenv("SMTP_HOST")
	SMTP_PORT := os.Getenv("SMTP_PORT")

	auth := smtp.PlainAuth("", SMTP_USER, SMTP_PASSWORD, SMTP_HOST)

	switch {
	case purpose == "OTPRESET":
		wd, wderr := os.Getwd()
		if wderr != nil {
			log.Fatal(wderr)
		}
		t, parsing_err := template.ParseFiles(wd + "/templates/otp_mail.html")
		if parsing_err != nil {
			fmt.Println("Error in parsing error ", parsing_err)
		}

		var body bytes.Buffer

		mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html;"
		body.Write([]byte(fmt.Sprintf("Subject: OTP to reset your password. \n%s\n\n", mimeHeaders)))

		t.Execute(&body, struct{ OTP string }{OTP: otp})

		err := smtp.SendMail(SMTP_HOST+":"+SMTP_PORT, auth, SMTP_USER, recipentList, body.Bytes())
		if err != nil {
			fmt.Println("ERR HERE > ", err)
			return
		}
		fmt.Println("MAIL SENT")
	}
}

// SendOTPToUser sends mail to the user
func SendOTPTOUser(email string, otp string) {
	recipent_list := []string{email}
	SendMailToUser(recipent_list, otp, "OTPRESET")
}

type Settings struct {
	DB_HOST       string
	DB_NAME       string
	DB_USER       string
	DB_PASSWORD   string
	DB_PORT       string
	SMTP_USER     string
	SMTP_PASSWORD string
	SMTP_HOST     string
	SMTP_PORT     string
	REDIS_HOST    string
}

func InitializeSettings() Settings {
	DB_USER := os.Getenv("DB_USER")
	DB_HOST := os.Getenv("DB_HOST")
	DB_NAME := os.Getenv("DB_NAME")
	DB_PASSWORD := os.Getenv("DB_PASSWORD")

	switch {
	case DB_HOST == "":
		fmt.Println("2 Environmet variable DB_HOST not set.")
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
