package main

import (
	"fmt"
	"hookfy/handlers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
		AllowCredentials: false,
	}))

	r.POST("/webhooks/:hash", handlers.CreateWebhook)

	r.GET("/webhooks/inbox", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Webhook inbox",
		})
	})

	fmt.Println("Server is running on port 8081")
	r.Run(":8081")
}
