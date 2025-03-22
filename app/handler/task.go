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
	"go.mongodb.org/mongo-driver/mongo/options"
)

// @Summary        Create Task
// @Description    Creates a new task
// @Tags           Task
// @Accept         json
// @Produce        json
// @Param          task body model.CreateTaskDTO true "Task Data"
// @Success        201 {object} Response
// @Router         /api/v1/tasks [post]
func CreateTask(c *fiber.Ctx) error {
	b := new(model.CreateTaskDTO)
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

	collUsers := db.GetDBCollection("users")
	count, err := collUsers.CountDocuments(c.Context(), bson.M{"_id": b.UserID})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(Response{
			Message: "Failed to check user_id uniqueness",
			Code:    http.StatusInternalServerError,
		})
	}
	if count == 0 {
		return c.Status(http.StatusBadRequest).JSON(Response{
			Message: "User not found",
			Code:    http.StatusBadRequest,
		})
	}
	if b.AssignedAt == nil {
		assignedAt := time.Now().UTC()
		b.AssignedAt = &assignedAt
	}
	if b.EstimatedPomodoros == nil {
		estimatedPomodoros := int16(1)
		b.EstimatedPomodoros = &estimatedPomodoros
	}

	task := model.CreateTaskDTO{
		UserID:             b.UserID,
		Title:              b.Title,
		Description:        b.Description,
		AssignedAt:         b.AssignedAt,
		Status:             model.TaskPending,
		EstimatedPomodoros: b.EstimatedPomodoros,
		CompletedPomodoros: b.CompletedPomodoros,
		CreatedAt:          time.Now().UTC(),
		UpdatedAt:          time.Now().UTC(),
	}

	coll := db.GetDBCollection("tasks")

	result, err := coll.InsertOne(c.Context(), task)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(Response{
			Message: "Failed to create task",
			Code:    http.StatusInternalServerError,
		})
	}

	return c.Status(http.StatusCreated).JSON(Response{
		Message: "Task created successfully",
		Code:    http.StatusCreated,
		Data:    result,
	})
}

// @Summary				Get Task by ID
// @Description		Retrieves a task from the database by ID
// @Tags					Task
// @Produce				json
// @Param					id path string true "Task ID"
// @Success				200 {object} Response
// @Router				/api/v1/tasks/{id} [get]
func GetTaskByID(c *fiber.Ctx) error {
	coll := db.GetDBCollection("tasks")

	id := c.Params("id")
	if id == "" {
		return c.Status(http.StatusBadRequest).JSON(Response{
			Message: "Invalid request",
			Code:    http.StatusBadRequest,
		})
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(Response{
			Message: "Invalid ID",
			Code:    http.StatusBadRequest,
		})
	}

	task := model.Task{}
	err = coll.FindOne(c.Context(), bson.M{"_id": objectID}).Decode(&task)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(Response{
			Message: "Task not found",
			Code:    http.StatusNotFound,
		})
	}

	return c.Status(http.StatusOK).JSON(Response{
		Message: "Task found",
		Code:    http.StatusOK,
		Data:    task,
	})
}

// @Summary				Update Task by ID
// @Description		Updates a task in the database by ID
// @Tags					Task
// @Accept				json
// @Produce				json
// @Param					id path string true "Task ID"
// @Param					task body model.UpdateTaskDTO true "Task Data"
// @Success				200 {object} Response
// @Router				/api/v1/tasks/{id} [put]
func UpdateTaskByID(c *fiber.Ctx) error {
	b := new(model.UpdateTaskDTO)
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

	id := c.Params("id")
	if id == "" {
		return c.Status(http.StatusBadRequest).JSON(Response{
			Message: "Invalid request",
			Code:    http.StatusBadRequest,
		})
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(Response{
			Message: "Invalid ID",
			Code:    http.StatusBadRequest,
		})
	}

	task := model.UpdateTaskDTO{
		Title:              b.Title,
		Description:        b.Description,
		AssignedAt:         b.AssignedAt,
		Status:             b.Status,
		EstimatedPomodoros: b.EstimatedPomodoros,
		CompletedPomodoros: b.CompletedPomodoros,
		UpdatedAt:          time.Now().UTC(),
	}

	coll := db.GetDBCollection("tasks")

	_, err = coll.UpdateOne(c.Context(), bson.M{"_id": objectID}, bson.M{"$set": task})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(Response{
			Message: "Failed to update task",
			Code:    http.StatusInternalServerError,
		})
	}

	return c.Status(http.StatusOK).JSON(Response{
		Message: "Task updated successfully",
		Code:    http.StatusOK,
	})
}

// @Summary				Delete Task by ID
// @Description		Deletes a task from the database by ID
// @Tags					Task
// @Produce				json
// @Param					id path string true "Task ID"
// @Success				200 {object} Response
// @Router				/api/v1/tasks/{id} [delete]
func DeleteTaskByID(c *fiber.Ctx) error {
	coll := db.GetDBCollection("tasks")

	id := c.Params("id")
	if id == "" {
		return c.Status(http.StatusBadRequest).JSON(Response{
			Message: "Invalid request",
			Code:    http.StatusBadRequest,
		})
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(Response{
			Message: "Invalid ID",
			Code:    http.StatusBadRequest,
		})
	}

	task := model.DeleteTaskDTO{
		Status:    model.TaskDeleted,
		DeletedAt: time.Now().UTC(),
	}

	_, err = coll.UpdateOne(c.Context(), bson.M{"_id": objectID}, bson.M{"$set": task})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(Response{
			Message: "Failed to delete task",
			Code:    http.StatusInternalServerError,
		})
	}

	return c.Status(http.StatusOK).JSON(Response{
		Message: "Task deleted successfully",
		Code:    http.StatusOK,
	})
}

// @Summary				Get Tasks by User ID
// @Description		Retrieves tasks from the database by User ID with optional filters and pagination
// @Tags					Task
// @Produce				json
// @Param					id path string true "User ID"
// @Param					status query string false "Task Status"
// @Param					title query string false "Task Title"
// @Param					start_date query string false "Start Date"
// @Param					end_date query string false "End Date"
// @Param					page query int false "Page number"
// @Param					limit query int false "Number of tasks per page"
// @Success				200 {object} Response
// @Router				/api/v1/tasks/user/{id} [get]
func GetTasksByUserID(c *fiber.Ctx) error {
	coll := db.GetDBCollection("tasks")

	id := c.Params("id")
	if id == "" {
		return c.Status(http.StatusBadRequest).JSON(Response{
			Message: "Invalid request",
			Code:    http.StatusBadRequest,
		})
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(Response{
			Message: "Invalid ID",
			Code:    http.StatusBadRequest,
		})
	}

	filter := bson.M{"user_id": objectID, "status": bson.M{"$ne": string(model.TaskDeleted)}}

	if status := c.Query("status"); status != "" {
		filter["status"] = status
	}
	if title := c.Query("title"); title != "" {
		filter["title"] = bson.M{"$regex": title, "$options": "i"}
	}
	if startDate := c.Query("start_date"); startDate != "" {
		start, err := time.Parse(time.RFC3339, startDate)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(Response{
				Message: "Invalid start date format",
				Code:    http.StatusBadRequest,
			})
		}
		filter["assigned_at"] = bson.M{"$gte": start}
	}
	if endDate := c.Query("end_date"); endDate != "" {
		end, err := time.Parse(time.RFC3339, endDate)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(Response{
				Message: "Invalid end date format",
				Code:    http.StatusBadRequest,
			})
		}
		if _, ok := filter["assigned_at"]; ok {
			filter["assigned_at"].(bson.M)["$lte"] = end
		} else {
			filter["assigned_at"] = bson.M{"$lte": end}
		}
	}

	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	skip := (page - 1) * limit

	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}).SetSkip(int64(skip)).SetLimit(int64(limit))

	cursor, err := coll.Find(c.Context(), filter, opts)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(Response{
			Message: "Failed to get tasks",
			Code:    http.StatusInternalServerError,
		})
	}

	var tasks []model.Task
	if err := cursor.All(c.Context(), &tasks); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(Response{
			Message: "Failed to get tasks",
			Code:    http.StatusInternalServerError,
		})
	}

	total, err := coll.CountDocuments(c.Context(), filter)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(Response{
			Message: "Failed to count tasks",
			Code:    http.StatusInternalServerError,
		})
	}

	return c.Status(http.StatusOK).JSON(Response{
		Message: "Tasks found",
		Code:    http.StatusOK,
		Data:    tasks,
		Total:   total,
	})
}
