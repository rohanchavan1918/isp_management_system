package main

import (
	"isp/controllers"
	"isp/models"
	"isp/routes"

	"github.com/joho/godotenv"
)

// @title ISP Management System API
// @version 1.0
// @description ISP management system.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /api/v1/
func main() {
	godotenv.Load()
	settings := controllers.InitializeSettings()
	_ = settings.DB_HOST
	models.ConnectDataBase()
	go controllers.InitialPlanCache()
	r := routes.SetupRouter()
	r.Run(":8080")
}
