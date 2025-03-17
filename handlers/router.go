package handlers

import (
	"github.com/gofiber/fiber/v2"
)

type Response struct {
	Message string      `json:"message"`
	Code    int         `json:"code"`
	Data    interface{} `json:"data,omitempty"`
}

func CreateRouter(r *fiber.App) {

	api := r.Group("/api")

	v1 := api.Group("/v1")

	v1.Get("/healthcheck", healthCheck)

}
