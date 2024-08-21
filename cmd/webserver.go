package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"wikitracker/internal"
)

func main() {
	group, err := internal.InitConsumerGroup()
	if err != nil {
		log.Fatalf("Failed to initialize consumer group: %v", err)
	}
	internal.StartConsuming(group)

	r := gin.Default()
	r.Use(cors.Default())

	r.LoadHTMLGlob("./internal/templates/*")
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// default page
	r.GET("/wikitrack")
}
