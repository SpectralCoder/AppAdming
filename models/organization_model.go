package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Orgnization struct {
	Id         primitive.ObjectID `json:"id,omitempty"`
	Name       string             `json:"name,omitempty" validate:"required"`
	Address    string             `json:"address,omitempty" validate:"required"`
	Phone      int                `json:"phone,omitempty" validate:"required"`
	Email      string             `json:"email,omitempty"`
	Created_at time.Time          `json:"created_at"`
}
