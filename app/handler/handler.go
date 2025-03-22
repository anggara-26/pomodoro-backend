package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type Response struct {
	Message string      `json:"message"`
	Code    int         `json:"code"`
	Data    interface{} `json:"data,omitempty"`
	Total   int64       `json:"total,omitempty"`
}

// HealthCheck checks the health of the server
//
//	@Summary        Health Check
//	@Description    Checks if the server is running
//	@Tags           Health
//	@Accept         json
//	@Produce        json
//	@Success        200 {object} Response
//	@Router         /api/v1/healthcheck [get]
func HealthCheck(c *fiber.Ctx) error {
	res := Response{
		Message: "Server is running",
		Code:    200,
	}

	return c.Status(http.StatusOK).JSON(res)
}
