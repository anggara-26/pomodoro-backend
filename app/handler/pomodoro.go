package handler

import (
	"net/http"
	"time"

	"github.com/anggara-26/pomodoro-backend.git/app/model"
	"github.com/anggara-26/pomodoro-backend.git/db"
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// @Summary        Start Pomodoro Session
// @Description    Starts a new pomodoro session
// @Tags           Pomodoro Session
// @Accept         json
// @Produce        json
// @Param          session body model.CreateSessionDTO true "Pomodoro Session Data"
// @Success        201 {object} Response
// @Router         /api/v1/sessions/start [post]
func StartPomodoroSession(c *fiber.Ctx) error {
	b := new(model.CreateSessionDTO)
	if err := c.BodyParser(b); err != nil {
		return c.Status(http.StatusBadRequest).JSON(Response{
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
	}

	validate := validator.New()
	if err := validate.Struct(b); err != nil {
		return c.Status(http.StatusBadRequest).JSON(Response{
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
	}

	collTasks := db.GetDBCollection("tasks")
	countTasks, err := collTasks.CountDocuments(c.Context(), bson.M{"_id": b.TaskID})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(Response{
			Message: "Failed to check task",
			Code:    http.StatusInternalServerError,
		})
	}
	if countTasks == 0 {
		return c.Status(http.StatusBadRequest).JSON(Response{
			Message: "Task not found",
			Code:    http.StatusBadRequest,
		})
	}

	coll := db.GetDBCollection("sessions")

	count, err := coll.CountDocuments(c.Context(), bson.M{"user_id": b.UserID, "status": model.SessionActive})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(Response{
			Message: "Failed to check active session",
			Code:    http.StatusInternalServerError,
		})
	}
	if count > 0 {
		_, err := coll.UpdateMany(c.Context(), bson.M{"user_id": b.UserID, "status": model.SessionActive}, bson.M{"$set": bson.M{"status": model.SessionBreak}})
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(Response{
				Message: "Failed to update active session status",
				Code:    http.StatusInternalServerError,
			})
		}
	}

	b.StartedAt = time.Now().UTC()
	b.Status = model.SessionActive

	result, err := coll.InsertOne(c.Context(), b)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(Response{
			Message: "Failed to create session",
			Code:    http.StatusInternalServerError,
		})
	}

	return c.Status(http.StatusCreated).JSON(Response{
		Message: "Session created successfully",
		Code:    http.StatusCreated,
		Data:    result.InsertedID,
	})
}

// @Summary        End Pomodoro Session
// @Description    Ends an active pomodoro session
// @Tags           Pomodoro Session
// @Accept         json
// @Produce        json
// @Param          id path string true "Session ID"
// @param          is_skip query bool false "Skip the session"
// @Success        200 {object} Response
// @Router         /api/v1/sessions/end/{id} [post]
func EndPomodoroSession(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(http.StatusBadRequest).JSON(Response{
			Message: "Session ID is required",
			Code:    http.StatusBadRequest,
		})
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(Response{
			Message: "Invalid session ID",
			Code:    http.StatusBadRequest,
		})
	}

	coll := db.GetDBCollection("sessions")

	count, err := coll.CountDocuments(c.Context(), bson.M{"_id": objectID, "status": model.SessionActive})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(Response{
			Message: "Failed to check active session",
			Code:    http.StatusInternalServerError,
		})
	}
	if count == 0 {
		return c.Status(http.StatusBadRequest).JSON(Response{
			Message: "Session not found or already ended",
			Code:    http.StatusBadRequest,
		})
	}

	b := model.EndSession{
		EndedAt: time.Now().UTC(),
		Status:  model.SessionCompleted,
	}

	if isSkipped := c.Query("is_skip"); isSkipped == "true" {
		b.Status = model.SessionSkipped
	}

	_, err = coll.UpdateOne(c.Context(), bson.M{"_id": objectID}, bson.M{"$set": b})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(Response{
			Message: "Failed to end session",
			Code:    http.StatusInternalServerError,
		})
	}

	return c.Status(http.StatusOK).JSON(Response{
		Message: "Session ended successfully",
		Code:    http.StatusOK,
	})
}
