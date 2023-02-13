package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type History struct {
	Id          primitive.ObjectID `json:"id,omitempty"`
	Due         int                `json:"Due,omitempty" validate:"required"`
	Paid        int                `json:"paid,omitempty" validate:"required"`
	Date        primitive.DateTime `json:"date,omitempty" validate:"required"`
	Customer_id primitive.ObjectID `json:"customer_id,omitempty" validate:"required"`
	Seller_id   primitive.ObjectID `json:"seller_id,omitempty" validate:"required"`
}
