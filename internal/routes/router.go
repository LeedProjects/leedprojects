package routes

import "github.com/gin-gonic/gin"

// LPRouter defines the application routes
func LPRouter() *gin.Engine {
	router := gin.Default()
	AddAuthRoutes(router)
	return router
}
