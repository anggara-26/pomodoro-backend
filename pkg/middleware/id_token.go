package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/anggara-26/pomodoro-backend.git/app/handler"
	"github.com/anggara-26/pomodoro-backend.git/platform/firebase"
	"github.com/gofiber/fiber/v2"
)

func FirebaseAuthMiddleware() fiber.Handler {
	authClient := firebase.InitFirebase()

	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(http.StatusUnauthorized).JSON(handler.Response{
				Message: "Authorization header is missing",
				Code:    http.StatusUnauthorized,
			})
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(http.StatusUnauthorized).JSON(handler.Response{
				Message: "Invalid authorization format",
				Code:    http.StatusUnauthorized,
			})
		}

		token, err := authClient.VerifyIDToken(context.Background(), parts[1])
		if err != nil {
			return c.Status(http.StatusUnauthorized).JSON(handler.Response{
				Message: "Invalid or expired token",
				Code:    http.StatusUnauthorized,
			})
		}

		c.Locals("firebase_uid", token.UID)
		return c.Next()
	}
}
