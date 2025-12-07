package main

import (
	"query-service/database"
	"query-service/events"
	"query-service/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	database.ConnectDB()

	// Start RabbitMQ Consumer (Writes)
	go events.StartConsumer()

	// Start HTTP Server (Reads)
	r := gin.Default()
	r.GET("/users/:email", handlers.GetUserHistory)

	r.Run(":8081") // Listen on port 8081
}
