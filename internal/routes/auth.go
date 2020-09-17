package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/ronyv89/leedprojects/internal/controllers"
)

// AddAuthRoutes adds authentication related routes to the router
func AddAuthRoutes(router *gin.Engine) {
	authRouter := router.Group("/auth")
	{
		authRouter.POST("/login", controllers.AuthLogin)
		authRouter.POST("/signup", controllers.AuthSignup)
	}
}
