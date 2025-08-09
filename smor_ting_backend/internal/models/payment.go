package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PaymentMethodType string

const (
	PaymentMethodMobileMoney PaymentMethodType = "mobile_money"
	PaymentMethodCard        PaymentMethodType = "card"
)

type LinkedPaymentMethod struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID    primitive.ObjectID `json:"user_id" bson:"user_id"`
	Type      PaymentMethodType  `json:"type" bson:"type"`
	Provider  string             `json:"provider" bson:"provider"`
	Msisdn    string             `json:"msisdn" bson:"msisdn"` // masked in responses
	IsDefault bool               `json:"is_default" bson:"is_default"`
	Status    string             `json:"status" bson:"status"` // active, pending, disabled
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}
