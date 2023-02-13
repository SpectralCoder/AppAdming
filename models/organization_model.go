package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Organization struct {
	Id       primitive.ObjectID `json:"id,omitempty"`
	Name     string             `json:"name,omitempty" validate:"required"`
	Thana    string             `json:"thana,omitempty" validate:"required"`
	District string             `json:"district,omitempty" validate:"required"`
	Phone    int                `json:"phone,omitempty" validate:"required"`
	Email    string             `json:"email,omitempty" validate:"required"`
}
