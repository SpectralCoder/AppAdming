package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type SellInfo struct {
	Id           primitive.ObjectID `json:"id,omitempty"`
	Products     []Product          `json:"products,omitempty" validate:"required"`
	Customer     Customer           `json:"customer,omitempty" validate:"required"`
	Organization Organization       `json:"seller,omitempty" validate:"required"`
	Amount       int                `json:"amount,omitempty" validate:"required"`
}
