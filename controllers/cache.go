package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"isp/models"
	"log"
	"os"
	"strconv"

	"github.com/go-redis/redis/v8"
	. "github.com/logrusorgru/aurora"
)

var ctx = context.Background()

// InitialPlanCache caches and sync all the plans from db to redis.
func InitialPlanCache() {
	// User Plans will be fetched by all users on dashboards and is common for all, thus it should be cached.
	fmt.Println(Bold(Cyan("[INFO] Started Caching ...")))
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST"),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	var plans []models.Plan
	db := models.DB

	result := db.Find(&plans)
	if result.Error != nil {
		log.Println("[ERR] ", result.Error)
	}
	rows, err := db.Model(&models.Plan{}).Rows()
	defer rows.Close()
	if err != nil {
		log.Println("[ERR] ", result.Error)
	}

	for rows.Next() {
		var plan models.Plan // ScanRows is a method of `gorm.DB`, it can be used to scan a row into a struct
		db.ScanRows(rows, &plan)
		var m = make(map[string]interface{})
		m["id"] = plan.ID
		m["name"] = plan.Name
		m["speed"] = plan.Speed
		m["duration"] = plan.Duration
		m["cost"] = plan.Cost
		m["notes"] = plan.Notes

		key := "plan:" + strconv.Itoa(plan.ID)
		doesKeyExist := CheckKeyExists(client, key, strconv.Itoa(plan.ID))
		if doesKeyExist == false {

			err := client.HMSet(ctx, key, m).Err()

			if err != nil {
				log.Println("[ERR] ", err)
			}
		}
	}

	// for models.DB.Next()
	log.Println(Bold(Cyan("[INFO] Initial plans cached.")))
}

func GetKeyInfo(client *redis.Client, key string) map[string]string {
	res, cerr := client.HGetAll(ctx, key).Result()
	if cerr != nil {
		fmt.Println(cerr)
	}
	return res
}

// CheckKeyExists checks if the key is present on redis
func CheckKeyExists(client *redis.Client, key string, id string) bool {
	_, err := client.HGet(ctx, key, id).Result()
	if err == redis.Nil {
		return false
	} else if err != nil {
		log.Println("caching err ")
	} else {
		return true
	}
	return false
}

func GetCachedPlans() []models.Plan {
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST"),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	cachedList := []models.Plan{}
	res, err := client.Keys(ctx, "plan:*").Result()
	if err == redis.Nil {
		// DB HIT
	}
	for _, v := range res {
		data := GetKeyInfo(client, v)
		// convert map to json
		jsonString, _ := json.Marshal(data)
		// convert json to struct
		s := models.Plan{}
		json.Unmarshal(jsonString, &s)
		cachedList = append(cachedList, s)
	}
	return cachedList
}
