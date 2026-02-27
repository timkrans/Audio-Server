package main

import(
	"log"

	"github.com/gin-gonic/gin"

	"audio-server/database"
	"audio-server/models"
	"audio-server/routes"
)

func main(){
	if err := database.Connect(); err != nil {
		log.Fatal(err)
	}
	database.DB.AutoMigrate(&models.Audio{})
	r := gin.Default()
	r.SetTrustedProxies(nil)
	routes.RegisterMovieRoutes(r)
	r.Run(":8080")
}