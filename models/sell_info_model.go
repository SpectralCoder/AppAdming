package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type SellInfo struct {
	Id              primitive.ObjectID `json:"id,omitempty"`
	Products        []*Product         `json:"products,omitempty" validate:"required"`
	Customer_id     primitive.ObjectID `json:"customer,omitempty" validate:"required"`
	Organization_id primitive.ObjectID `json:"seller,omitempty" validate:"required"`
	Amount          int                `json:"amount,omitempty" validate:"required"`
}
