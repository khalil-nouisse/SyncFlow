package handlers

import (
	"context"
	"net/http"
	"query-service/database"
	"query-service/models"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

// GetUserHistory handles GET /users/:email
func GetUserHistory(c *gin.Context) {
	email := c.Param("email")

	// 1. Define Context with Timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 2. Find the document
	var user models.User
	err := database.UserCollection.FindOne(ctx, bson.M{"_id": email}).Decode(&user)

	if err != nil {
		// If generic error (includes 'not found')
		c.JSON(http.StatusNotFound, gin.H{"error": "User or history not found"})
		return
	}

	// 3. Return JSON
	c.JSON(http.StatusOK, user)
}
