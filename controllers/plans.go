package controllers

import (
	"context"
	"fmt"
	"html/template"
	"isp/models"
	"log"
	"net/http"
	"strconv"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	. "github.com/logrusorgru/aurora"
	"gorm.io/gorm"
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
		CacheNewPlan(sanitized_new_plan)
		c.JSON(http.StatusCreated, gin.H{"id": sanitized_new_plan.ID, "name": sanitized_new_plan.Name})
		fmt.Println("caching new plan")

		return
	}
}

func CacheNewPlan(plan models.Plan) {
	// ScanRows is a method of `gorm.DB`, it can be used to scan a row into a struct
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	var ctx = context.Background()

	var m = make(map[string]interface{})
	m["id"] = plan.ID
	m["name"] = plan.Name
	m["speed"] = plan.Speed
	m["duration"] = plan.Duration
	m["cost"] = plan.Cost
	m["notes"] = plan.Notes

	key := "plan:" + strconv.Itoa(plan.ID)
	if !CheckKeyExists(client, key, strconv.Itoa(plan.ID)) {
		err := client.HSet(ctx, key, m)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(Bold(Cyan("[INFO] New plan cached.")))
		log.Println(Bold(Cyan("[INFO] New plan cached.")))
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
		// Get all cached plans

		models.DB.Find(&all_plans)
		c.JSON(http.StatusOK, &all_plans)
	}
}

func GetAllCachedPlans(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": GetCachedPlans()})
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

func UpdatePlan(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	user_email, _ := claims["email"]
	var User models.User
	var Plan models.Plan
	type UpdatePlanInput struct {
		gorm.Model
		Cost     int
		Name     string
		Speed    string
		Duration int    // Number of days
		Notes    string //Additoinal string	}
	}

	plan_id := c.Param("id")

	if err := models.DB.Where("email = ? AND role=1", user_email).First(&User).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	// Main logic here.
	if err := models.DB.Where("id = ?", plan_id).First(&Plan).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	var input UpdatePlanInput

	if err := c.ShouldBindJSON(&input); err != nil {
		fmt.Println("Error in bind json", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println("plan is ", Plan)
	fmt.Println("Updates to be done ", input)
	models.DB.Model(&Plan).Updates(input)
	c.JSON(http.StatusOK, gin.H{"data": &Plan})

}

func GetPlan(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	user_email, _ := claims["email"]
	var User models.User
	var Plan models.Plan

	plan_id := c.Param("id")

	if err := models.DB.Where("email = ? AND role=1", user_email).First(&User).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	// Main logic here.
	if err := models.DB.Where("id = ?", plan_id).First(&Plan).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": &Plan})

}
