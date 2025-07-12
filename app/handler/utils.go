package handler

import (
	"net/http"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// ErrorResponse represents an error response
type ErrorResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Field   string `json:"field,omitempty"`
}

// ValidationError represents validation error details
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// HandleValidationErrors handles validator errors and returns formatted response
func HandleValidationErrors(c *fiber.Ctx, err error) error {
	var errors []ValidationError

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldError := range validationErrors {
			var message string
			switch fieldError.Tag() {
			case "required":
				message = "This field is required"
			case "email":
				message = "Must be a valid email address"
			case "min":
				message = "Value is too small"
			case "max":
				message = "Value is too large"
			default:
				message = "Invalid value"
			}

			errors = append(errors, ValidationError{
				Field:   fieldError.Field(),
				Message: message,
			})
		}
	}

	return c.Status(http.StatusBadRequest).JSON(Response{
		Message: "Validation failed",
		Code:    http.StatusBadRequest,
		Data:    errors,
	})
}

// HandleMongoError handles MongoDB errors and returns appropriate response
func HandleMongoError(c *fiber.Ctx, err error) error {
	if err == mongo.ErrNoDocuments {
		return c.Status(http.StatusNotFound).JSON(Response{
			Message: "Resource not found",
			Code:    http.StatusNotFound,
		})
	}

	return c.Status(http.StatusInternalServerError).JSON(Response{
		Message: "Database error occurred",
		Code:    http.StatusInternalServerError,
	})
}

// ValidateObjectID validates and converts string to ObjectID
func ValidateObjectID(id string) (primitive.ObjectID, error) {
	return primitive.ObjectIDFromHex(id)
}

// SuccessResponse returns a success response
func SuccessResponse(c *fiber.Ctx, message string, data interface{}) error {
	return c.Status(http.StatusOK).JSON(Response{
		Message: message,
		Code:    http.StatusOK,
		Data:    data,
	})
}

// CreatedResponse returns a created response
func CreatedResponse(c *fiber.Ctx, message string, data interface{}) error {
	return c.Status(http.StatusCreated).JSON(Response{
		Message: message,
		Code:    http.StatusCreated,
		Data:    data,
	})
}

// ErrorResponseWithStatus returns an error response with custom status
func ErrorResponseWithStatus(c *fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(Response{
		Message: message,
		Code:    status,
	})
}

// BadRequestResponse returns a bad request response
func BadRequestResponse(c *fiber.Ctx, message string) error {
	return ErrorResponseWithStatus(c, http.StatusBadRequest, message)
}

// NotFoundResponse returns a not found response
func NotFoundResponse(c *fiber.Ctx, message string) error {
	return ErrorResponseWithStatus(c, http.StatusNotFound, message)
}

// InternalServerErrorResponse returns an internal server error response
func InternalServerErrorResponse(c *fiber.Ctx, message string) error {
	return ErrorResponseWithStatus(c, http.StatusInternalServerError, message)
}

// UnauthorizedResponse returns an unauthorized response
func UnauthorizedResponse(c *fiber.Ctx, message string) error {
	return ErrorResponseWithStatus(c, http.StatusUnauthorized, message)
}
