package main

import (
	"context"
	"hookfy/config"
	"hookfy/handlers"
	"hookfy/worker"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	port, exists := os.LookupEnv("PORT")
	if !exists {
		log.Fatalln("PORT variable not set")
	}

	config.Connect()
	worker.StartDeleteExpiredWorker(config.DB, context.Background())

	r := gin.Default()

	r.Static("/static", "web/static")
	r.LoadHTMLGlob("web/templates/*")

	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
		AllowCredentials: false,
	}))

	r.POST("/webhooks/:hash", handlers.CreateWebhook)
	r.GET("/webhooks/inbox", handlers.GetInbox)
	r.GET("/webhooks/:id", handlers.GetWebhook)

	r.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", nil)
	})

	r.Run(":" + port)
}
