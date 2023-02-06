package models

import (
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// #1- omitempty => to send the filled data
// #2- primitive.ObjectId => MongoDB.ObjectId
// #3- uuid(v4) with MongoDB

type Book struct {
	ID          uuid.UUID          `json:"id,omitempty"`
	CreatedDate primitive.DateTime `json:"createdDate,omitempty"`
	UpdatedDate primitive.DateTime `json:"updatedDate"`
	Title       string             `json:"title,omitempty"`
	Author      string             `json:"author,omitempty"`
	Quantity    int                `json:"quantity,omitempty"`
}
