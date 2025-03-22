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

// @Summary        Create User
// @Description    Creates a new user in the database
// @Tags           User
// @Accept         json
// @Produce        json
// @Param          user body model.CreateUserDTO true "User Data"
// @Success        201 {object} Response
// @Router         /api/v1/users [post]
func CreateUser(c *fiber.Ctx) error {
	b := new(model.CreateUserDTO)
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

	coll := db.GetDBCollection("users")

	count, err := coll.CountDocuments(c.Context(), bson.M{"firebase_uid": b.FirebaseUID})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(Response{
			Message: "Failed to check firebase_uid uniqueness",
			Code:    http.StatusInternalServerError,
		})
	}
	if count > 0 {
		return c.Status(http.StatusBadRequest).JSON(Response{
			Message: "FirebaseUID already exists",
			Code:    http.StatusBadRequest,
		})
	}

	b.CreatedAt = time.Now().UTC()

	result, err := coll.InsertOne(c.Context(), b)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(Response{
			Message: "Failed to create user",
			Code:    http.StatusInternalServerError,
		})
	}

	return c.Status(http.StatusCreated).JSON(Response{
		Message: "User created successfully",
		Code:    http.StatusCreated,
		Data:    result,
	})
}

// @Summary        Get User by ID
// @Description    Retrieves a user from the database by ID
// @Tags           User
// @Produce        json
// @Param          id path string true "User ID"
// @Success        200 {object} Response
// @Router         /api/v1/users/{id} [get]
func GetUserByID(c *fiber.Ctx) error {
	coll := db.GetDBCollection("users")

	id := c.Params("id")
	if id == "" {
		return c.Status(http.StatusBadRequest).JSON(Response{
			Message: "Invalid request",
			Code:    http.StatusBadRequest,
		})
	}

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(Response{
			Message: "Invalid ID",
			Code:    http.StatusBadRequest,
		})
	}

	user := model.User{}

	err = coll.FindOne(c.Context(), bson.M{"_id": objectId}).Decode(&user)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(Response{
			Message: "User not found",
			Code:    http.StatusNotFound,
		})
	}

	return c.Status(http.StatusOK).JSON(Response{
		Message: "User found",
		Code:    http.StatusOK,
		Data:    user,
	})
}

// @Summary        Update User by ID
// @Description    Updates a user in the database by ID
// @Tags           User
// @Accept         json
// @Produce        json
// @Param          id path string true "User ID"
// @Param          user body model.UpdateUserDTO true "User Data"
// @Success        200 {object} Response
// @Router         /api/v1/users/{id} [put]
func UpdateUserByID(c *fiber.Ctx) error {
	coll := db.GetDBCollection("users")

	id := c.Params("id")
	if id == "" {
		return c.Status(http.StatusBadRequest).JSON(Response{
			Message: "Invalid request",
			Code:    http.StatusBadRequest,
		})
	}

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(Response{
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
	}

	b := new(model.UpdateUserDTO)
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

	_, err = coll.UpdateOne(c.Context(), bson.M{"_id": objectId}, bson.M{"$set": b})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(Response{
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
	}

	return c.Status(http.StatusOK).JSON(Response{
		Message: "User updated successfully",
		Code:    http.StatusOK,
	})
}
