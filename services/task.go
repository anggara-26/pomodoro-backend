package services

import "go.mongodb.org/mongo-driver/mongo"

type Task struct {
	ID                 string `json:"id,omitempty" bson:"_id,omitempty"`
	Title              string `json:"title,omitempty" bson:"title,omitempty"`
	Description        string `json:"description,omitempty" bson:"description,omitempty"`
	AssignedAt         string `json:"assigned_at,omitempty" bson:"assigned_at,omitempty"`
	Status             string `json:"status,omitempty" bson:"status,omitempty"`
	EstimatedPomodoros int    `json:"estimated_pomodoros,omitempty" bson:"estimated_pomodoros,omitempty"`
	CompletedPomodoros int    `json:"completed_pomodoros,omitempty" bson:"completed_pomodoros,omitempty"`
	CreatedAt          string `json:"created_at,omitempty" bson:"created_at,omitempty"`
}

var client *mongo.Client

func New(mongo *mongo.Client) Task {
	client = mongo

	return Task{}
}
