package handlers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func healthCheck(c *fiber.Ctx) error {
	res := Response{
		Message: "Server is running",
		Code:    200,
		Data:    "Pomodoro Backend",
	}

	return c.Status(http.StatusOK).JSON(res)
}
