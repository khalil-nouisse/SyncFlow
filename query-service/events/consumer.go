package events

import (
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
	EventType string          `json:"event_type"`
	Payload   json.RawMessage `json:"payload"`
}

// ProductEventPayload matches the payload sent by producer
type ProductEventPayload struct {
	Description string  `json:"p_desc"`
	Quantity    float64 `json:"qte"`
}

func StartConsumer() {
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

	// Queue name must match producer: 'product_events'
	q, _ := ch.QueueDeclare("product_events", true, false, false, false, nil)
	msgs, _ := ch.Consume(q.Name, "", true, false, false, false, nil)

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			var wrapper EventWrapper
			if err := json.Unmarshal(d.Body, &wrapper); err != nil {
				log.Printf("❌ Error parsing wrapper: %v", err)
				continue
			}

			switch wrapper.EventType {
			case "PRODUCT_CREATED":
				var payload ProductEventPayload
				if err := json.Unmarshal(wrapper.Payload, &payload); err != nil {
					log.Printf("❌ Error parsing product payload: %v", err)
					continue
				}
				handleProductInsert(payload)

			case "ORDER_CREATED":
				var payload models.Order
				if err := json.Unmarshal(wrapper.Payload, &payload); err != nil {
					log.Printf("❌ Error parsing order creation payload: %v", err)
					continue
				}
				handleOrderInsert(payload)

			case "ORDER_UPDATED":
				var payload models.Order
				if err := json.Unmarshal(wrapper.Payload, &payload); err != nil {
					log.Printf("❌ Error parsing order update payload: %v", err)
					continue
				}
				handleOrderUpdate(payload)

			case "ORDER_DELETED":
				var payload models.Order
				if err := json.Unmarshal(wrapper.Payload, &payload); err != nil {
					log.Printf("❌ Error parsing order delete payload: %v", err)
					continue
				}
				handleOrderDelete(payload)

			default:
				log.Printf("⚠️ Unknown Event Type: %s", wrapper.EventType)
			}
		}
	}()

	log.Printf("🎧 Waiting for events. Press CTRL+C to exit")
	<-forever
}

// --- PRODUCT HANDLERS ---
func handleProductInsert(p ProductEventPayload) {
	product := models.Product{
		Description: p.Description,
		Quantity:    p.Quantity,
	}
	_, err := database.ProductCollection.InsertOne(nil, product)
	if err != nil {
		log.Printf("❌ Product Insert Failed: %v", err)
	} else {
		log.Printf("✅ Inserted Product: %s", p.Description)
	}
}

// --- ORDER HANDLERS ---

func handleOrderInsert(order models.Order) {
	// Force ID to be the same if we want to query by it later easily,
	// or let Mongo generate one. Since we receive 'id_commande', let's set _id to it.
	// However, models.Order already has `bson:"_id,omitempty"`.
	// If id_commande is populated in the struct, Mongo driver might use it if we map it to _id.
	// The struct definition I made earlier:
	// ID string `bson:"_id,omitempty" json:"id_commande"`
	// So JSON unmarshal puts `id_commande` into `ID`.
	// BSON insertion will use `ID` as `_id`. Perfect.

	_, err := database.OrderCollection.InsertOne(nil, order)
	if err != nil {
		log.Printf("❌ Order Insert Failed: %v", err)
	} else {
		log.Printf("✅ Inserted Order: %s", order.ID)
	}
}

func handleOrderUpdate(order models.Order) {
	filter := bson.M{"_id": order.ID}
	update := bson.M{"$set": order}
	opts := options.Update().SetUpsert(true) // Upsert to be safe? Or just update.

	_, err := database.OrderCollection.UpdateOne(nil, filter, update, opts)
	if err != nil {
		log.Printf("❌ Order Update Failed: %v", err)
	} else {
		log.Printf("🔄 Updated Order: %s", order.ID)
	}
}

func handleOrderDelete(order models.Order) {
	filter := bson.M{"_id": order.ID}
	_, err := database.OrderCollection.DeleteOne(nil, filter)
	if err != nil {
		log.Printf("❌ Order Delete Failed: %v", err)
	} else {
		log.Printf("🗑️ Deleted Order: %s", order.ID)
	}
}
