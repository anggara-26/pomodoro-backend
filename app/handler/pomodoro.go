package handler

import (
	"log"
	"net/http"
	"time"

	"github.com/anggara-26/pomodoro-backend.git/app/model"
	"github.com/anggara-26/pomodoro-backend.git/platform/db"
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
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

	// Validate user exists
	collUsers := db.GetDBCollection("users")
	countUsers, err := collUsers.CountDocuments(c.Context(), bson.M{"_id": b.UserID})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(Response{
			Message: "Failed to check user",
			Code:    http.StatusInternalServerError,
		})
	}
	if countUsers == 0 {
		return c.Status(http.StatusBadRequest).JSON(Response{
			Message: "User not found",
			Code:    http.StatusBadRequest,
		})
	}

	// Validate task exists if provided
	if b.TaskID != nil {
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

	// Get the session first to check if it exists and get task info
	var session model.Session
	err = coll.FindOne(c.Context(), bson.M{"_id": objectID, "status": model.SessionActive}).Decode(&session)
	if err != nil {
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

	// If session was completed (not skipped) and has a task, update completed pomodoros
	if b.Status == model.SessionCompleted && session.TaskID != primitive.NilObjectID && session.Type == model.Focus {
		tasksColl := db.GetDBCollection("tasks")
		_, err = tasksColl.UpdateOne(
			c.Context(),
			bson.M{"_id": session.TaskID},
			bson.M{"$inc": bson.M{"completed_pomodoros": 1}},
		)
		if err != nil {
			log.Printf("Warning: Failed to update task pomodoro count: %v", err)
		}

		// Check if task should be marked as completed
		var task model.Task
		err = tasksColl.FindOne(c.Context(), bson.M{"_id": session.TaskID}).Decode(&task)
		if err == nil && task.CompletedPomodoros+1 >= task.EstimatedPomodoros {
			tasksColl.UpdateOne(
				c.Context(),
				bson.M{"_id": session.TaskID},
				bson.M{"$set": bson.M{"status": model.TaskCompleted, "updated_at": time.Now().UTC()}},
			)
		}
	}

	return c.Status(http.StatusOK).JSON(Response{
		Message: "Session ended successfully",
		Code:    http.StatusOK,
	})
}

// @Summary        Get Active Sessions
// @Description    Gets all active sessions for a user
// @Tags           Pomodoro Session
// @Produce        json
// @Param          user_id query string true "User ID"
// @Success        200 {object} Response
// @Router         /api/v1/sessions/active [get]
func GetActiveSessions(c *fiber.Ctx) error {
	userID := c.Query("user_id")
	if userID == "" {
		return BadRequestResponse(c, "User ID is required")
	}

	userObjectID, err := ValidateObjectID(userID)
	if err != nil {
		return BadRequestResponse(c, "Invalid user ID")
	}

	coll := db.GetDBCollection("sessions")
	filter := bson.M{
		"user_id": userObjectID,
		"status":  model.SessionActive,
	}

	cursor, err := coll.Find(c.Context(), filter)
	if err != nil {
		return InternalServerErrorResponse(c, "Failed to get active sessions")
	}

	var sessions []model.Session
	if err := cursor.All(c.Context(), &sessions); err != nil {
		return InternalServerErrorResponse(c, "Failed to decode sessions")
	}

	return SuccessResponse(c, "Active sessions retrieved successfully", sessions)
}

// @Summary        Get Session History
// @Description    Gets session history for a user with pagination and filtering
// @Tags           Pomodoro Session
// @Produce        json
// @Param          user_id query string true "User ID"
// @Param          type query string false "Session Type (focus, short_break, long_break)"
// @Param          status query string false "Session Status (active, completed, skipped)"
// @Param          start_date query string false "Start Date (YYYY-MM-DD)"
// @Param          end_date query string false "End Date (YYYY-MM-DD)"
// @Param          page query int false "Page number"
// @Param          limit query int false "Number of sessions per page"
// @Success        200 {object} Response
// @Router         /api/v1/sessions/history [get]
func GetSessionHistory(c *fiber.Ctx) error {
	userID := c.Query("user_id")
	if userID == "" {
		return BadRequestResponse(c, "User ID is required")
	}

	userObjectID, err := ValidateObjectID(userID)
	if err != nil {
		return BadRequestResponse(c, "Invalid user ID")
	}

	filter := bson.M{"user_id": userObjectID}

	// Add optional filters
	if sessionType := c.Query("type"); sessionType != "" {
		filter["type"] = sessionType
	}
	if status := c.Query("status"); status != "" {
		filter["status"] = status
	}
	if startDate := c.Query("start_date"); startDate != "" {
		start, err := time.Parse("2006-01-02", startDate)
		if err != nil {
			return BadRequestResponse(c, "Invalid start date format. Use YYYY-MM-DD")
		}
		filter["started_at"] = bson.M{"$gte": start}
	}
	if endDate := c.Query("end_date"); endDate != "" {
		end, err := time.Parse("2006-01-02", endDate)
		if err != nil {
			return BadRequestResponse(c, "Invalid end date format. Use YYYY-MM-DD")
		}
		// Add one day to include the end date
		endPlusOne := end.Add(24 * time.Hour)
		if _, ok := filter["started_at"]; ok {
			filter["started_at"].(bson.M)["$lt"] = endPlusOne
		} else {
			filter["started_at"] = bson.M{"$lt": endPlusOne}
		}
	}

	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 20)
	skip := (page - 1) * limit

	coll := db.GetDBCollection("sessions")
	opts := options.Find().
		SetSort(bson.D{{Key: "started_at", Value: -1}}).
		SetSkip(int64(skip)).
		SetLimit(int64(limit))

	cursor, err := coll.Find(c.Context(), filter, opts)
	if err != nil {
		return InternalServerErrorResponse(c, "Failed to get session history")
	}

	var sessions []model.Session
	if err := cursor.All(c.Context(), &sessions); err != nil {
		return InternalServerErrorResponse(c, "Failed to decode sessions")
	}

	total, err := coll.CountDocuments(c.Context(), filter)
	if err != nil {
		return InternalServerErrorResponse(c, "Failed to count sessions")
	}

	return c.Status(http.StatusOK).JSON(Response{
		Message: "Session history retrieved successfully",
		Code:    http.StatusOK,
		Data:    sessions,
		Total:   total,
	})
}
