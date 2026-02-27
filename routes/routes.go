package routes

import(

	"github.com/gin-gonic/gin"
	"net/http"
)

func RegisterMovieRoutes(r *gin.Engine) {
	//only use * for testing
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Header("Access-Control-Allow-Methods", "GET, OPTIONS")
		c.Header("Cache-Control", "no-cache")
	})
		//adding a health check for later testing
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
}
