package main

import (
	"isp/controllers"
	"isp/models"
	"isp/routes"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	settings := controllers.InitializeSettings()
	_ = settings.DB_HOST
	models.ConnectDataBase()
	go controllers.InitialPlanCache()
	r := routes.SetupRouter()
	r.Run(":8080")
}
