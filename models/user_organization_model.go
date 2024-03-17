package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type UserOrganization struct {
	Id           primitive.ObjectID `json:"id,omitempty"`
	Organization primitive.ObjectID `json:"organization,omitempty" validate:"required"`
	User         primitive.ObjectID `json:"user,omitempty" validate:"required"`
	Role         string             `json:"role,omitempty" validate:"required"`
	Status       string             `json:"status,omitempty" validate:"required"`
}

type UserOrganizationRequest struct {
	Id     primitive.ObjectID `json:"id,omitempty"`
	Role   string             `json:"role,omitempty" validate:"required"`
	Status string             `json:"status,omitempty" validate:"required"`
}
