package controllers

import (
	"errors"
	"fmt"
	"isp/models"
	"net/http"
	"strconv"
	"text/template"
	"time"

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

type ForgotPasswordInput struct {
	Email string `binding:"required"`
}

type VerifyOTPInput struct {
	Email string `binding:"required"`
	OTP   string `binding:"required"`
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

// @Summary Signup/Register/Add Users
// @Description Signup/Register/Add users
// @Router /api/v1/signup [post]
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
		Role:        2,
		Password:    encryptedPassword,
		IsActive:    true,
	}

	models.DB.Create(&sanitizeduser)
	c.JSON(http.StatusCreated, gin.H{"msg": "User created successfully"})

}

// @Summary /api/v1/auth/whoami returns the basic details (id, email) of the logged user.
// @Description returns the ID, Email of the currently loggedin user.
// @Router /api/v1/auth/whoami [get]
// @Accept json
// @Produce json
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
			"role":  User.Role,
		})
	}
}

// @Summary /api/v1/auth/reset_password ResetPassword allows you to reset your password.
// @Description returns the ID, Email of the currently loggedin user.
// @Router /api/v1/auth/reset_password [post]
// @Accept json
// @Produce json
// @Param user body ResetPasswordInput true "User Data"
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

// Creates a OTP for the user.
func ForgotPassword(c *gin.Context) {
	var forgot_password_input ForgotPasswordInput
	var User models.User
	var ForgetPassword models.ForgotPassword

	if err := c.BindJSON(&forgot_password_input); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request."})
		return
	}
	if err := models.DB.Where("email = ?", template.HTMLEscapeString(forgot_password_input.Email)).First(&User).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User does not exist."})
		return
	}
	ForgetPassword.User = User
	ForgetPassword.UserID = int(User.ID)
	otp := GetOtp()
	ForgetPassword.OTP = int(otp)
	ForgetPassword.Created_at = time.Now()
	ForgetPassword.Expired_at = time.Now().Add(time.Minute * 5)
	models.DB.Save(&ForgetPassword)

	go SendOTPTOUser(forgot_password_input.Email, strconv.Itoa(int(otp)))
	c.JSON(http.StatusOK, gin.H{"message": "OTP Generated successfully."})
}

func VerifyOTP(c *gin.Context) {
	var verify_otp VerifyOTPInput
	// var User models.User
	var Forgot_Password models.ForgotPassword
	if err := c.BindJSON(&verify_otp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request."})
		return
	}
	// Check if empty string has been passed in email
	if verify_otp.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request."})
		return
	}
	// Check if the len of otp is 6
	if len(verify_otp.OTP) != 6 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request."})
		return
	}
	// type FindOtpResults struct {
	var (
		FPID       string
		UID        string
		Email      string
		OTP        int
		Expired_at time.Time
	)

	var affected_rows int64
	// var Find_otp_results FindOtpResults
	models.DB.Raw(`select count(fp.id) from forgot_passwords fp join users u on fp.user_id=u.id where u.email=? and fp.otp=? ;`, verify_otp.Email, verify_otp.OTP).Scan(&affected_rows)
	if affected_rows == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid OTP or Email"})
		return
	}
	// var Find_otp_results FindOtpResults
	rows, err := models.DB.Raw(`select fp.id,u.id,u.email,fp.otp,fp.expired_at from forgot_passwords fp join users u on fp.user_id=u.id where u.email=? and fp.otp=? ;`, verify_otp.Email, verify_otp.OTP).Rows()
	defer rows.Close()
	if err != nil {
		return
	}
	for rows.Next() {
		rows.Scan(&FPID, &UID, &Email, &OTP, &Expired_at)
	}
	// Check if the token is valid.
	if time.Now().Before(Expired_at) == false {
		c.JSON(http.StatusBadRequest, gin.H{"error": "OTP has been expired"})
		return
	}
	// if all condition matches, set the token.
	reset_secured_token := GenerateSecureToken(16)
	if err := models.DB.Model(&Forgot_Password).Where("id = ?", &FPID).Update("token", reset_secured_token).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"messsage": "Some error occoured"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": reset_secured_token})
}

// Update password
// api/v1/update_forgotten_password/assdasd12qe2cqec2
// {
// 	password
// 	new_password
// }

func UpdateForgotenPassword(c *gin.Context) {
	token := c.Param("token")
	// var token_count int64
	var user models.User
	var forgot_password models.ForgotPassword
	type UpdatePasswordInput struct {
		Password        string `json:"password" binding:"required";`
		ConfirmPassword string `json:"confirm_password" binding:"required";`
	}

	var update_password_input UpdatePasswordInput

	// First check if the token exists in database
	res := models.DB.First(&forgot_password, "token=?", token)
	if res.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "This link is invalid"})
		return
	}

	if time.Now().Before(forgot_password.Expired_at) == false {
		c.JSON(http.StatusNotFound, gin.H{"error": "This link has been expired"})
		return
	}

	if err := c.BindJSON(&update_password_input); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request."})
		return
	}
	// Check password matches the password policy.
	is_password_strong, _ := IsPasswordStrong(template.HTMLEscapeString(update_password_input.Password))
	if is_password_strong != true {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password is not strong enough."})
		return
	}
	if update_password_input.Password != update_password_input.ConfirmPassword {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Both password doesnot match."})
		return
	}

	test_res := models.DB.First(&user, forgot_password.UserID)
	if test_res.Error != nil {
		fmt.Println("err", res.Error)
	}
	new_encrypted_password, error := HashPassword(template.HTMLEscapeString(update_password_input.Password))
	if error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Some error occoured."})
	}

	user.Password = new_encrypted_password
	models.DB.Save(&user)

	forgot_password.Expired_at = time.Now()
	models.DB.Save(&forgot_password)

	c.JSON(http.StatusOK, gin.H{"message": "Password updated successfully."})
}
