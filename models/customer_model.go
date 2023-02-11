package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Customer struct {
	Id       primitive.ObjectID `json:"id,omitempty"`
	Name     string             `json:"name,omitempty" validate:"required"`
	Father   string             `json:"father,omitempty" validate:"required"`
	Home     string             `json:"home,omitempty" validate:"required"`
	Village  string             `json:"village,omitempty" validate:"required"`
	Thana    string             `json:"thana,omitempty" validate:"required"`
	District string             `json:"district,omitempty" validate:"required"`
	Paid     int                `json:"paid,omitempty"`
	Due      int                `json:"due,omitempty"`
}
