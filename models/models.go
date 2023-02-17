package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// #1- omitempty => to send the filled data
// #2- primitive.ObjectId => MongoDB.ObjectId
// #3- uuid(v4) with MongoDB

type Book struct {
	ID          string             `json:"_id,omitempty" bson:"_id,omitempty"`
	CreatedDate primitive.DateTime `json:"createddate,omitempty"`
	UpdatedDate primitive.DateTime `json:"updateddate"`
	Title       string             `json:"title,omitempty"`
	Author      string             `json:"author,omitempty"`
	Quantity    int                `json:"quantity,omitempty"`
}
