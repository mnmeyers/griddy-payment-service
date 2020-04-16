package payment

import (
	"github.com/stripe/stripe-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// I'm assuming client will send amount in the form of xx.xx
type PostBody struct {
	CustomerId string `json:"account_id"`
	Amount     string `json:"amount"` // instructions required Amount to be string
}

type DTO struct {
	PaymentId string `json:"id"`
	Amount    string `json:"amount"` // instructions required Amount to be string
	Status    string `json:"status"`
}

type DAO struct {
	ID         primitive.ObjectID         `bson:"_id,omitempty"`
	Amount     int64                      `bson:"amount"`
	CustomerId string                     `bson:"customerId"`
	PaymentId  string                     `bson:"paymentId"`
	Status     stripe.PaymentIntentStatus `bson:"status"`
}
