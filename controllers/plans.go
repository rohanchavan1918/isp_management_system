package controllers

import (
	"context"
	"fmt"
	"html/template"
	"isp/models"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	. "github.com/logrusorgru/aurora"
)

type UpdatePlanInput struct {
	ID       int    `json:"id,string"`
	Name     string `json:"name"`
	Speed    string `json:"speed"`
	Duration int    `json:"duration"` // Number of days
	Cost     int    `json:"cost"`
	Notes    string `json:"notes"` //Additoinal string
	IsActive bool   `json:"isactive"`
}

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

func DeletePlanFromCache(key string) {
	fmt.Println("DELETING cache ", key)
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST"),
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	var ctx = context.Background()
	res := client.Del(ctx, key)
	if res.Val() == 0 {
		log.Println("Not Deleted from cache ", key)
	} else if res.Val() == 1 {
		log.Println("Deleted from cache ", key)
	}
}

func CacheNewPlan(plan models.Plan) {
	// ScanRows is a method of `gorm.DB`, it can be used to scan a row into a struct
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST"),
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
		log.Println(Bold(Cyan("[INFO] New plan cached.")))
	}

}

// GetAllPlans returns the list of plans
// all?source=0 -- Returns data from cache
// all?source=1 -- Returns data from DB
func GetAllPlans(c *gin.Context) {

	user_source := c.Query("source")

	claims := jwt.ExtractClaims(c)
	user_email, _ := claims["email"]
	var User models.User
	if user_source == "0" {
		if err := models.DB.Where("email = ?", user_email).First(&User).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		} else {
			// Get all cached plans
			c.JSON(http.StatusOK, gin.H{"data": GetCachedPlans()})
			return
		}
	} else if user_source == "1" {
		all_plans := []models.Plan{}

		if err := models.DB.Where("email = ?", user_email).First(&User).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		} else {
			models.DB.Find(&all_plans)
			c.JSON(http.StatusOK, &all_plans)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "Not found."})
	return
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
		go DeletePlanFromCache("plan:" + plan_id)
		c.JSON(http.StatusOK, gin.H{"message": "Plan deleted successfully"})
	}
}

func UpdatePlanInCache(key string, plan models.Plan) {
	DeletePlanFromCache(key)
	CacheNewPlan(plan)
}

func UpdatePlan(c *gin.Context) {
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

	var input UpdatePlanInput

	if err := c.ShouldBindJSON(&input); err != nil {
		fmt.Println("Error in bind json", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	models.DB.Model(&Plan).Updates(input)
	// Update the values incache ( Delete existing entry in cache, add new entry in cache  )
	UpdatePlanInCache("plan:"+strconv.Itoa(int(input.ID)), Plan)

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

func AddUserToPlan(c *gin.Context) {

	type UserPlanInput struct {
		UserID int `json:"user_id"`
		PlanID int `json:"plan_id"`
	}
	claims := jwt.ExtractClaims(c)
	user_email, _ := claims["email"]
	var User models.User
	var TargetUser models.User
	var Plan models.Plan
	var Input UserPlanInput
	var NewUserPlan models.UserPlans

	log.Println("user email ", user_email)
	if err := models.DB.Where("email = ? AND role=1", user_email).First(&User).Error; err != nil {
		log.Println("err asd a", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	log.Println("Got user ", User.FirstName)
	if err := c.ShouldBindJSON(&Input); err != nil {
		fmt.Println("Error in bind json", err)
		c.JSON(http.StatusBadRequest, gin.H{"error here": err.Error()})
		return
	}

	if err := models.DB.Where("id = ? AND role=2", Input.UserID).First(&TargetUser).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found."})
		return
	}

	if serr := models.DB.First(&Plan, int(Input.PlanID)).Error; serr != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": serr})
	}

	now := time.Now()
	then := now.AddDate(0, 0, Plan.Duration)

	NewUserPlan.PlanId = Plan.ID
	NewUserPlan.UserId = int(TargetUser.ID)
	NewUserPlan.IsActive = true
	NewUserPlan.ValidTill = then
	NewUserPlan.Created_at = now

	res := models.DB.Create(&NewUserPlan)
	if res.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Some error occoured."})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": NewUserPlan})

}
