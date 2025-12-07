package events

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"query-service/database"
	"query-service/models"

	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// EventWrapper is the outer shell of every message
type EventWrapper struct {
	EventType string          `json:"event_type"` // e.g., "ORDER_CREATED", "ORDER_UPDATED"
	Payload   json.RawMessage `json:"payload"`    // Delays parsing until we know the type
}

// OrderEventPayload is the data inside the wrapper
type OrderEventPayload struct {
	UserEmail string       `json:"user_email"`
	UserName  string       `json:"user_name"`
	Order     models.Order `json:"order"`
}

func StartConsumer() {
	// ... (Connection logic remains the same as before) ...
	// [Copy lines 30-70 from previous consumer.go here: Get Env, Connect, Channel, QueueDeclare]
	// Here is the shortened version for the loop:

	url := os.Getenv("RABBITMQ_URL")
	if url == "" {
		url = "amqp://guest:guest@localhost:5672/"
	}

	var conn *amqp.Connection
	var err error
	maxRetries := 15
	
	for i := 0; i < maxRetries; i++ {
		log.Printf("Attempting to connect to RabbitMQ (Attempt %d/%d)...", i+1, maxRetries)
		conn, err = amqp.Dial(url)
		if err == nil {
			log.Println("✅ Connected to RabbitMQ successfully")
			break
		}
		log.Printf("❌ Failed to connect to RabbitMQ: %v. Retrying in 5 seconds...", err)
		time.Sleep(5 * time.Second)
	}

	if err != nil {
		log.Fatalf("Could not connect to RabbitMQ after %d attempts: %v", maxRetries, err)
	}
	defer conn.Close()
	ch, _ := conn.Channel()
	defer ch.Close()
	q, _ := ch.QueueDeclare("order_created_queue", true, false, false, false, nil)
	msgs, _ := ch.Consume(q.Name, "", true, false, false, false, nil)

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			// 1. Unmarshal only the Wrapper to check the Type
			var wrapper EventWrapper
			if err := json.Unmarshal(d.Body, &wrapper); err != nil {
				log.Printf("❌ Error parsing wrapper: %v", err)
				continue
			}

			// 2. Parse the inner data
			var payload OrderEventPayload
			if err := json.Unmarshal(wrapper.Payload, &payload); err != nil {
				log.Printf("❌ Error parsing payload: %v", err)
				continue
			}

			// 3. Switch based on Event Type
			switch wrapper.EventType {
			case "ORDER_CREATED":
				handleInsert(payload)
			case "ORDER_UPDATED":
				handleUpdate(payload)
			case "ORDER_DELETED":
				handleDelete(payload)
			default:
				log.Printf("⚠️ Unknown Event Type: %s", wrapper.EventType)
			}
		}
	}()

	log.Printf("🎧 Waiting for events. Press CTRL+C to exit")
	<-forever
}

// --- CRUD HANDLERS ---

// 1. INSERT (Create)
func handleInsert(p OrderEventPayload) {
	filter := bson.M{"_id": p.UserEmail}
	update := bson.M{
		"$setOnInsert": bson.M{"name": p.UserName},
		"$push":        bson.M{"orders": p.Order},
	}
	opts := options.Update().SetUpsert(true)
	_, err := database.UserCollection.UpdateOne(context.TODO(), filter, update, opts)
	if err != nil {
		log.Printf("❌ Insert Failed: %v", err)
	} else {
		log.Printf("✅ Inserted Order %d for %s", p.Order.OrderID, p.UserEmail)
	}
}

// 2. UPDATE (Modify an existing order status/total)
func handleUpdate(p OrderEventPayload) {
	// Filter: Find user by Email AND the specific order inside the array
	filter := bson.M{
		"_id":             p.UserEmail,
		"orders.order_id": p.Order.OrderID,
	}

	// Update: Use the positional operator '$' to update exactly that order
	update := bson.M{
		"$set": bson.M{
			"orders.$.status":       p.Order.Status,      // Update Status
			"orders.$.total_amount": p.Order.TotalAmount, // Update Amount
		},
	}

	_, err := database.UserCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Printf("❌ Update Failed: %v", err)
	} else {
		log.Printf("🔄 Updated Order %d for %s", p.Order.OrderID, p.UserEmail)
	}
}

// 3. DELETE (Remove an order)
func handleDelete(p OrderEventPayload) {
	filter := bson.M{"_id": p.UserEmail}

	// Update: Use $pull to remove the item from the array where order_id matches
	update := bson.M{
		"$pull": bson.M{
			"orders": bson.M{"order_id": p.Order.OrderID},
		},
	}

	_, err := database.UserCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Printf("❌ Delete Failed: %v", err)
	} else {
		log.Printf("🗑️ Deleted Order %d from %s", p.Order.OrderID, p.UserEmail)
	}
}
