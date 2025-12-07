package models

type OrderItem struct {
	ProductID int `bson:"product_id" json:"productId"`
	Qty       int `bson:"qty" json:"qty"`
}

type Order struct {
	OrderID     int         `bson:"order_id" json:"orderId"`
	Items       []OrderItem `bson:"items" json:"items"`
	TotalAmount float64     `bson:"total_amount" json:"totalAmount"`
	Status      string      `bson:"status" json:"status"`
}

type User struct {
	UserID    int     `bson:"user_id" json:"userId"`
	FirstName string  `bson:"first_name" json:"firstName"`
	LastName  string  `bson:"last_name" json:"lastName"`
	Address   string  `bson:"address" json:"address"`
	Orders    []Order `bson:"orders" json:"orders"`
}
