package main

import(
	"log"

	"github.com/gin-gonic/gin"

	"audio-server/database"
	"audio-server/models"
)

func main(){
	if err := database.Connect(); err != nil {
		log.Fatal(err)
	}
	database.DB.AutoMigrate(&models.Audio{})
	r := gin.Default()
	r.SetTrustedProxies(nil)
	r.Run(":8080")
}