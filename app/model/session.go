package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Session struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	UserID    primitive.ObjectID `json:"user_id" bson:"user_id"`
	TaskID    primitive.ObjectID `json:"task_id" bson:"task_id"`
	StartedAt time.Time          `json:"started_at" bson:"started_at"`
	EndedAt   time.Time          `json:"ended_at" bson:"ended_at"`
	Duration  int16              `json:"duration" bson:"duration"`
	Type      SessionType        `json:"type" bson:"type"`
	Status    SessionStatus      `json:"status" bson:"status"`
}

type CreateSessionDTO struct {
	UserID    primitive.ObjectID  `json:"user_id" bson:"user_id" validate:"required"`
	TaskID    *primitive.ObjectID `json:"task_id,omitempty" bson:"task_id,omitempty"`
	StartedAt time.Time           `json:"started_at" bson:"started_at"`
	EndedAt   time.Time           `json:"ended_at" bson:"ended_at"`
	Duration  int16               `json:"duration" bson:"duration" validate:"required"`
	Type      SessionType         `json:"type" bson:"type" validate:"required"`
	Status    SessionStatus       `json:"status" bson:"status"`
}

type EndSession struct {
	EndedAt time.Time     `json:"ended_at" bson:"ended_at"`
	Status  SessionStatus `json:"status" bson:"status"`
}

type SessionType string

const (
	Focus      SessionType = "focus"
	ShortBreak SessionType = "short_break"
	LongBreak  SessionType = "long_break"
)

type SessionStatus string

const (
	SessionActive    SessionStatus = "active"
	SessionBreak     SessionStatus = "break"
	SessionSkipped   SessionStatus = "skipped"
	SessionCompleted SessionStatus = "completed"
)
