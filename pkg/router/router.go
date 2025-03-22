package router

import (
	"github.com/anggara-26/pomodoro-backend.git/app/handler"
	_ "github.com/anggara-26/pomodoro-backend.git/docs/v1"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
)

func CreateRouter(r *fiber.App) {
	r.Use(logger.New())
	r.Use(recover.New())
	r.Use(cors.New())

	r.Get("/swagger/*", swagger.HandlerDefault)

	api := r.Group("/api")

	v1 := api.Group("/v1")
	v1.Get("/healthcheck", handler.HealthCheck)

	users := v1.Group("/users")
	users.Post("/", handler.CreateUser)
	users.Get("/:id", handler.GetUserByID)
	users.Put("/:id", handler.UpdateUserByID)

	tasks := v1.Group("/tasks")
	tasks.Post("/", handler.CreateTask)
	tasks.Get("/user/:id", handler.GetTasksByUserID)
	tasks.Get("/:id", handler.GetTaskByID)
	tasks.Put("/:id", handler.UpdateTaskByID)
	tasks.Delete("/:id", handler.DeleteTaskByID)

	sessions := v1.Group("/sessions")
	sessions.Post("/start", handler.StartPomodoroSession)
	sessions.Post("/end/:id", handler.EndPomodoroSession)

	// stats := v1.Group("/stats", middleware.AuthMiddleware)
	// stats.Get("/daily", getDailyStats)
	// stats.Get("/weekly", getWeeklyStats)
	// stats.Get("/monthly", getMonthlyStats)
}
