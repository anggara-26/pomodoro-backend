package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	FirebaseUID string             `json:"firebase_uid" bson:"firebase_uid" validate:"required"`
	Email       string             `json:"email" bson:"email" validate:"required"`
	Name        string             `json:"name" bson:"name" validate:"required"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
}

type CreateUserDTO struct {
	FirebaseUID string    `json:"firebase_uid" bson:"firebase_uid" validate:"required"`
	Email       string    `json:"email" bson:"email" validate:"required"`
	Name        string    `json:"name" bson:"name" validate:"required"`
	CreatedAt   time.Time `json:"created_at" bson:"created_at"`
}

type UpdateUserDTO struct {
	Name string `json:"name" bson:"name" validate:"required"`
}
