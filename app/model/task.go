package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Task struct {
	ID                 primitive.ObjectID `json:"id" bson:"_id"`
	UserID             primitive.ObjectID `json:"user_id" bson:"user_id"`
	Title              string             `json:"title" bson:"title"`
	Description        *string            `json:"description" bson:"description"`
	AssignedAt         time.Time          `json:"assigned_at" bson:"assigned_at"`
	Status             TaskStatus         `json:"status" bson:"status"`
	EstimatedPomodoros int16              `json:"estimated_pomodoros" bson:"estimated_pomodoros" validate:"min=1"`
	CompletedPomodoros int16              `json:"completed_pomodoros" bson:"completed_pomodoros"`
	CreatedAt          time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt          time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	DeletedAt          *time.Time         `json:"deleted_at,omitempty" bson:"deleted_at,omitempty"`
}

type CreateTaskDTO struct {
	UserID             primitive.ObjectID `json:"user_id" bson:"user_id" validate:"required"`
	Title              string             `json:"title" bson:"title" validate:"required"`
	Description        *string            `json:"description,omitempty" bson:"description,omitempty"`
	AssignedAt         *time.Time         `json:"assigned_at,omitempty" bson:"assigned_at,omitempty"`
	Status             TaskStatus         `json:"status" bson:"status"`
	EstimatedPomodoros *int16             `json:"estimated_pomodoros" bson:"estimated_pomodoros" validate:"omitempty,min=1"`
	CompletedPomodoros int16              `json:"completed_pomodoros" bson:"completed_pomodoros"`
	CreatedAt          time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt          time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	DeletedAt          time.Time          `json:"deleted_at,omitempty" bson:"deleted_at,omitempty"`
}

type UpdateTaskDTO struct {
	Title              *string     `json:"title,omitempty" bson:"title,omitempty"`
	Description        *string     `json:"description,omitempty" bson:"description,omitempty"`
	AssignedAt         *time.Time  `json:"assigned_at,omitempty" bson:"assigned_at,omitempty"`
	Status             *TaskStatus `json:"status,omitempty" bson:"status,omitempty"`
	EstimatedPomodoros *int16      `json:"estimated_pomodoros,omitempty" bson:"estimated_pomodoros,omitempty" validate:"omitempty,min=1"`
	CompletedPomodoros *int16      `json:"completed_pomodoros,omitempty" bson:"completed_pomodoros,omitempty"`
	UpdatedAt          time.Time   `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

type DeleteTaskDTO struct {
	Status    TaskStatus `json:"status" bson:"status"`
	DeletedAt time.Time  `json:"deleted_at" bson:"deleted_at"`
}

type TaskStatus string

const (
	TaskPending    TaskStatus = "pending"
	TaskInProgress TaskStatus = "in_progress"
	TaskCompleted  TaskStatus = "completed"
	TaskDeleted    TaskStatus = "deleted"
)
