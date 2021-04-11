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
	r := routes.SetupRouter()
	r.Run(":8000")
}
