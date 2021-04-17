package routes

import (
	"isp/controllers"
	"isp/middlewares"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// PublicEndpoints returns signup signin
func PublicEndpoints(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {

	r.POST("/signup", controllers.Signup)
	r.GET("/signup", func(c *gin.Context) {
		c.JSON(401, gin.H{
			"error": "Method Not Allowed",
		})
	})

	r.POST("/signin", authMiddleware.LoginHandler)
	r.GET("/signin", func(c *gin.Context) {
		c.JSON(401, gin.H{
			"error": "Method Not Allowed",
		})
	})

	// Forgot password
	r.POST("/forgot_password", controllers.ForgotPassword)
	r.GET("/forgot_password", func(c *gin.Context) {
		c.JSON(401, gin.H{
			"error": "Method Not Allowed",
		})
	})

	// Forgot password
	r.POST("/verify_otp", controllers.VerifyOTP)
	r.GET("/verify_otp", func(c *gin.Context) {
		c.JSON(401, gin.H{
			"error": "Method Not Allowed",
		})
	})

	// Update forgoten Password
	r.POST("/update_forgoten_password/:token", controllers.UpdateForgotenPassword)
	r.GET("/update_forgoten_password", func(c *gin.Context) {
		c.JSON(401, gin.H{
			"error": "Method Not Allowed",
		})
	})

}

// SetupRouter creates a router with default
func SetupRouter() *gin.Engine {
	gin.ForceConsoleColor()
	r := gin.Default()
	authMiddleware, _ := middlewares.GetAuthMiddleware()

	v1 := r.Group("/api/v1/")
	PublicEndpoints(v1.Group(""), authMiddleware)
	auth := r.Group("/api/v1/auth")
	auth.GET("/refresh_token", authMiddleware.RefreshHandler)
	auth.Use(authMiddleware.MiddlewareFunc())
	{
		auth.GET("/whoami", controllers.GetIDFromEmail)
		auth.POST("/reset_password", controllers.ResetPassword)
		auth.POST("/plans/add", controllers.AddPlan)
		auth.GET("/plans/all", controllers.GetAllPlans)
		auth.GET("/plans/cached_all", controllers.GetAllCachedPlans)
		auth.GET("/plan/:id", controllers.GetPlan)
		auth.PATCH("/plan/:id", controllers.UpdatePlan)
		auth.DELETE("/plan/:id", controllers.DeletePlan)
	}

	return r
}
