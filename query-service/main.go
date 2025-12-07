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
	r.GET("/products", handlers.GetAllProducts)
	r.GET("/orders", handlers.GetAllOrders)
	r.GET("/orders/:id", handlers.GetOrder)

	r.Run(":8081") // Listen on port 8081
}
