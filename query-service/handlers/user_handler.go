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

// GetProduct handles GET /product/:id
// Note: Originally user/:email, assuming we want to fetch by ID or some field.
// Let's stick to user request of "do the same thing in mongodb" which implies getting the data.
// Since the producer sends `p_desc` and `qte`, we might not have a clean ID unless we use the mongo ID.
// For now, let's just make it fetch all or fetch by some criteria.
// Given the previous code fetched by "email" which was the ID, let's fetch by ID if possible, or just list all.
// Actually, the previous code was `GetUserHistory`.
// Let's make `GetAllProducts` for simplicity as `getAll` is in producer.

func GetAllProducts(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := database.ProductCollection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch products"})
		return
	}
	defer cursor.Close(ctx)

	var products []models.Product
	if err = cursor.All(ctx, &products); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse products"})
		return
	}

	c.JSON(http.StatusOK, products)
}
