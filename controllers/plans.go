package controllers

import (
	"fmt"
	"html/template"
	"isp/models"
	"net/http"
	"strconv"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

func AddPlan(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	user_email, _ := claims["email"]
	var User models.User

	var new_plan models.Plan
	var sanitized_new_plan models.Plan

	if err := models.DB.Where("email = ? AND role=1", user_email).First(&User).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	} else {
		// c.JSON(200, gin.H{
		// 	"id":    User.ID,
		// 	"email": User.Email,
		// 	"role":  User.Role,
		// })

		if err := c.BindJSON(&new_plan); err != nil {
			fmt.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request."})
			return
		}
		sanitized_new_plan.Name = template.HTMLEscapeString(new_plan.Name)
		sanitized_new_plan.Speed = template.HTMLEscapeString(new_plan.Speed)
		sanitized_new_plan.Duration, _ = strconv.Atoi(template.HTMLEscapeString(strconv.Itoa(new_plan.Duration)))
		sanitized_new_plan.Cost, _ = strconv.Atoi(template.HTMLEscapeString(strconv.Itoa(new_plan.Cost)))
		sanitized_new_plan.Notes = template.HTMLEscapeString(new_plan.Notes)

		res := models.DB.Create(&sanitized_new_plan)
		if res.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Some error occoured."})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"id": sanitized_new_plan.ID, "name": sanitized_new_plan.Name})

	}
}

func GetAllPlans(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	user_email, _ := claims["email"]
	var User models.User

	all_plans := []models.Plan{}

	if err := models.DB.Where("email = ? AND role=1", user_email).First(&User).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	} else {
		models.DB.Find(&all_plans)
		c.JSON(http.StatusOK, &all_plans)
	}
}

func DeletePlan(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	user_email, _ := claims["email"]
	var User models.User
	var plan models.Plan
	plan_id := c.Param("id")

	if err := models.DB.Where("email = ? AND role=1", user_email).First(&User).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	} else {

		models.DB.Delete(&plan, plan_id)
		c.JSON(http.StatusOK, gin.H{"message": "Plan deleted successfully"})
	}
}
