package controllers

import (
	"errors"
	"fmt"
	"isp/models"
	"net/http"
	"text/template"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type TempUser struct {
	FirstName   string `binding:"required"`
	LastName    string `binding:"required"`
	Email       string `binding:"required"`
	MobileNo    string `binding:"required"`
	DateOfBirth string `binding:"required"`
	Gender      string `binding:"required"`
	Password    string `binding:"required"`
}

type ResetPasswordInput struct {
	CurrentPassword    string `binding:"required"`
	NewPassword        string `binding:"required"`
	ConfirmNewPassword string `binding:"required"`
}

// Err returns a singular type err, which returns the err with line
type Err struct {
	message string
}

// Errors is a list of err
type Errors struct {
	list []Err
}

func DoesUserExist(email string) bool {
	var users []models.User
	err := models.DB.Where("email=?", email).First(&users).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false
		}
	}
	return true
}

// Signup returns the registeration info
func Signup(c *gin.Context) {
	var tempuser TempUser

	if err := c.BindJSON(&tempuser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ispasswordstrong, _ := IsPasswordStrong(tempuser.Password)
	if ispasswordstrong == false {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password is not strong."})
		return
	}

	if DoesUserExist(tempuser.Email) {

		c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists."})
		return
	}

	if tempuser.Gender == "F" || tempuser.Gender == "M" {
		// Gender
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid value in Gender"})
		return
	}

	encryptedPassword, error := HashPassword(template.HTMLEscapeString(tempuser.Password))
	if error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Some error occoured."})
	}

	sanitizeduser := models.User{
		FirstName:   template.HTMLEscapeString(tempuser.FirstName),
		LastName:    template.HTMLEscapeString(tempuser.LastName),
		Email:       template.HTMLEscapeString(tempuser.Email),
		MobileNo:    template.HTMLEscapeString(tempuser.MobileNo),
		DateOfBirth: template.HTMLEscapeString(tempuser.DateOfBirth),
		Gender:      template.HTMLEscapeString(tempuser.Gender),
		Password:    encryptedPassword,
		IsActive:    true,
	}

	models.DB.Create(&sanitizeduser)
	c.JSON(http.StatusCreated, gin.H{"msg": "User created successfully"})

}

// GetIDFromEmail returns the id and other info from the email
func GetIDFromEmail(c *gin.Context) {
	// db := c.MustGet("db").(*gorm.DB)
	claims := jwt.ExtractClaims(c)
	user_email, _ := claims["email"]
	var User models.User

	if err := models.DB.Where("email = ?", user_email).First(&User).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		c.JSON(200, gin.H{
			"id":    User.ID,
			"email": User.Email,
		})
	}
}

// ResetPassword, reset the password for authenticated user.
func ResetPassword(c *gin.Context) {
	/*
		Authenticates the user and then reset the password.
		1. Check if the entered currect password is correct.
		2. check if newpassword and currentnewpassword are equal.
		3. check if the new password matches the password policy (Strong password)
		4. Update password.
	*/
	claims := jwt.ExtractClaims(c)
	user_email, _ := claims["email"]
	var User models.User
	var reset_password_input ResetPasswordInput

	if err := c.BindJSON(&reset_password_input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// First check if the new password, and confirm password are same.
	if reset_password_input.NewPassword != reset_password_input.ConfirmNewPassword {
		c.JSON(http.StatusBadRequest, gin.H{"error": "new password doesnot match confirm password."})
		return
	}
	isNewPwdStrong, _ := IsPasswordStrong(reset_password_input.NewPassword)
	if isNewPwdStrong == false {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password is not strong."})
		return
	}

	if err := models.DB.Where("email = ?", user_email).First(&User).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	}

	err := bcrypt.CompareHashAndPassword([]byte(User.Password), []byte(reset_password_input.CurrentPassword))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "current password doesnot match."})
		return
	}

	// ALL CHECKS PASSED, password can be updated with the new password.
	User.Password, err = HashPassword(template.HTMLEscapeString(reset_password_input.NewPassword))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Some Error Occoured."})
		return
	}

	models.DB.Save(&User)

	c.JSON(201, gin.H{"message": "Password updated successfully."})

}
