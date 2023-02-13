package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type History struct {
	Id           primitive.ObjectID `json:"id,omitempty"`
	Due          int                `json:"Due,omitempty" validate:"required"`
	Paid         int                `json:"paid,omitempty" validate:"required"`
	Date         primitive.DateTime `json:"date,omitempty" validate:"required"`
	Customer     Customer           `json:"customer,omitempty" validate:"required"`
	Organization Organization       `json:"organization,omitempty" validate:"required"`
}
