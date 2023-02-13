package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Product struct {
	Id          primitive.ObjectID `json:"id,omitempty"`
	Model       string             `json:"model,omitempty" validate:"required"`
	Price       float64            `json:"price,omitempty" validate:"required"`
	Cost        float64            `json:"cost,omitempty" validate:"required"`
	Description string             `json:"description,omitempty"`
	Category    string             `json:"category,omitempty" validate:"required"`
	ImageURL    string             `json:"image_url,omitempty"`
	Stock       int                `json:"stock,omitempty" validate:"required"`
	Brand       int                `json:"brand,omitempty" validate:"required"`
	Seller_id   primitive.ObjectID `json:"organization,omitempty" validate:"required"`
}
