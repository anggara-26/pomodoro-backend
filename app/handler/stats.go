package handler

import (
	"net/http"
	"time"

	"github.com/anggara-26/pomodoro-backend.git/app/model"
	"github.com/anggara-26/pomodoro-backend.git/platform/db"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DailyStats struct {
	Date               string `json:"date"`
	CompletedPomodoros int    `json:"completed_pomodoros"`
	CompletedTasks     int    `json:"completed_tasks"`
	TotalTime          int    `json:"total_time_minutes"`
}

type WeeklyStats struct {
	WeekStart          string       `json:"week_start"`
	WeekEnd            string       `json:"week_end"`
	CompletedPomodoros int          `json:"completed_pomodoros"`
	CompletedTasks     int          `json:"completed_tasks"`
	TotalTime          int          `json:"total_time_minutes"`
	DailyBreakdown     []DailyStats `json:"daily_breakdown"`
}

type MonthlyStats struct {
	Month              string        `json:"month"`
	Year               int           `json:"year"`
	CompletedPomodoros int           `json:"completed_pomodoros"`
	CompletedTasks     int           `json:"completed_tasks"`
	TotalTime          int           `json:"total_time_minutes"`
	WeeklyBreakdown    []WeeklyStats `json:"weekly_breakdown"`
}

// @Summary        Get Daily Statistics
// @Description    Gets daily statistics for a user
// @Tags           Statistics
// @Accept         json
// @Produce        json
// @Param          user_id query string true "User ID"
// @Param          date query string false "Date (YYYY-MM-DD format)"
// @Success        200 {object} Response
// @Security       BearerAuth
// @Router         /api/v1/stats/daily [get]
func GetDailyStats(c *fiber.Ctx) error {
	userID := c.Query("user_id")
	if userID == "" {
		return c.Status(http.StatusBadRequest).JSON(Response{
			Message: "User ID is required",
			Code:    http.StatusBadRequest,
		})
	}

	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(Response{
			Message: "Invalid user ID",
			Code:    http.StatusBadRequest,
		})
	}

	// Parse date or use today
	dateStr := c.Query("date")
	var targetDate time.Time
	if dateStr != "" {
		targetDate, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(Response{
				Message: "Invalid date format. Use YYYY-MM-DD",
				Code:    http.StatusBadRequest,
			})
		}
	} else {
		targetDate = time.Now().UTC()
	}

	startOfDay := time.Date(targetDate.Year(), targetDate.Month(), targetDate.Day(), 0, 0, 0, 0, time.UTC)
	endOfDay := startOfDay.Add(24 * time.Hour)

	// Get completed sessions for the day
	sessionsColl := db.GetDBCollection("sessions")
	sessionsFilter := bson.M{
		"user_id":    userObjectID,
		"status":     model.SessionCompleted,
		"type":       model.Focus,
		"started_at": bson.M{"$gte": startOfDay, "$lt": endOfDay},
	}

	cursor, err := sessionsColl.Find(c.Context(), sessionsFilter)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(Response{
			Message: "Failed to get sessions",
			Code:    http.StatusInternalServerError,
		})
	}

	var sessions []model.Session
	if err := cursor.All(c.Context(), &sessions); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(Response{
			Message: "Failed to decode sessions",
			Code:    http.StatusInternalServerError,
		})
	}

	// Get completed tasks for the day
	tasksColl := db.GetDBCollection("tasks")
	tasksFilter := bson.M{
		"user_id":    userObjectID,
		"status":     model.TaskCompleted,
		"updated_at": bson.M{"$gte": startOfDay, "$lt": endOfDay},
	}

	completedTasksCount, err := tasksColl.CountDocuments(c.Context(), tasksFilter)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(Response{
			Message: "Failed to count completed tasks",
			Code:    http.StatusInternalServerError,
		})
	}

	// Calculate total time
	totalTime := 0
	for _, session := range sessions {
		totalTime += int(session.Duration)
	}

	stats := DailyStats{
		Date:               targetDate.Format("2006-01-02"),
		CompletedPomodoros: len(sessions),
		CompletedTasks:     int(completedTasksCount),
		TotalTime:          totalTime,
	}

	return c.Status(http.StatusOK).JSON(Response{
		Message: "Daily statistics retrieved successfully",
		Code:    http.StatusOK,
		Data:    stats,
	})
}

// @Summary        Get Weekly Statistics
// @Description    Gets weekly statistics for a user
// @Tags           Statistics
// @Accept         json
// @Produce        json
// @Param          user_id query string true "User ID"
// @Param          week_start query string false "Week start date (YYYY-MM-DD format)"
// @Success        200 {object} Response
// @Security       BearerAuth
// @Router         /api/v1/stats/weekly [get]
func GetWeeklyStats(c *fiber.Ctx) error {
	userID := c.Query("user_id")
	if userID == "" {
		return c.Status(http.StatusBadRequest).JSON(Response{
			Message: "User ID is required",
			Code:    http.StatusBadRequest,
		})
	}

	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(Response{
			Message: "Invalid user ID",
			Code:    http.StatusBadRequest,
		})
	}

	// Parse week start or use current week
	weekStartStr := c.Query("week_start")
	var weekStart time.Time
	if weekStartStr != "" {
		weekStart, err = time.Parse("2006-01-02", weekStartStr)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(Response{
				Message: "Invalid date format. Use YYYY-MM-DD",
				Code:    http.StatusBadRequest,
			})
		}
	} else {
		now := time.Now().UTC()
		weekday := int(now.Weekday())
		if weekday == 0 {
			weekday = 7 // Sunday = 7
		}
		weekStart = now.AddDate(0, 0, -weekday+1) // Monday
	}

	weekStart = time.Date(weekStart.Year(), weekStart.Month(), weekStart.Day(), 0, 0, 0, 0, time.UTC)
	weekEnd := weekStart.Add(7 * 24 * time.Hour)

	// Get daily breakdown
	var dailyStats []DailyStats
	totalPomodoros := 0
	totalTasks := 0
	totalTime := 0

	for i := 0; i < 7; i++ {
		currentDay := weekStart.Add(time.Duration(i) * 24 * time.Hour)
		nextDay := currentDay.Add(24 * time.Hour)

		// Get sessions for this day
		sessionsColl := db.GetDBCollection("sessions")
		sessionsFilter := bson.M{
			"user_id":    userObjectID,
			"status":     model.SessionCompleted,
			"type":       model.Focus,
			"started_at": bson.M{"$gte": currentDay, "$lt": nextDay},
		}

		cursor, err := sessionsColl.Find(c.Context(), sessionsFilter)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(Response{
				Message: "Failed to get sessions",
				Code:    http.StatusInternalServerError,
			})
		}

		var sessions []model.Session
		cursor.All(c.Context(), &sessions)

		// Get completed tasks for this day
		tasksColl := db.GetDBCollection("tasks")
		tasksFilter := bson.M{
			"user_id":    userObjectID,
			"status":     model.TaskCompleted,
			"updated_at": bson.M{"$gte": currentDay, "$lt": nextDay},
		}

		completedTasksCount, _ := tasksColl.CountDocuments(c.Context(), tasksFilter)

		// Calculate daily total time
		dayTotalTime := 0
		for _, session := range sessions {
			dayTotalTime += int(session.Duration)
		}

		dailyStats = append(dailyStats, DailyStats{
			Date:               currentDay.Format("2006-01-02"),
			CompletedPomodoros: len(sessions),
			CompletedTasks:     int(completedTasksCount),
			TotalTime:          dayTotalTime,
		})

		totalPomodoros += len(sessions)
		totalTasks += int(completedTasksCount)
		totalTime += dayTotalTime
	}

	stats := WeeklyStats{
		WeekStart:          weekStart.Format("2006-01-02"),
		WeekEnd:            weekEnd.Add(-time.Second).Format("2006-01-02"),
		CompletedPomodoros: totalPomodoros,
		CompletedTasks:     totalTasks,
		TotalTime:          totalTime,
		DailyBreakdown:     dailyStats,
	}

	return c.Status(http.StatusOK).JSON(Response{
		Message: "Weekly statistics retrieved successfully",
		Code:    http.StatusOK,
		Data:    stats,
	})
}

// @Summary        Get Monthly Statistics
// @Description    Gets monthly statistics for a user
// @Tags           Statistics
// @Accept         json
// @Produce        json
// @Param          user_id query string true "User ID"
// @Param          month query int false "Month (1-12)"
// @Param          year query int false "Year"
// @Success        200 {object} Response
// @Security       BearerAuth
// @Router         /api/v1/stats/monthly [get]
func GetMonthlyStats(c *fiber.Ctx) error {
	userID := c.Query("user_id")
	if userID == "" {
		return c.Status(http.StatusBadRequest).JSON(Response{
			Message: "User ID is required",
			Code:    http.StatusBadRequest,
		})
	}

	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(Response{
			Message: "Invalid user ID",
			Code:    http.StatusBadRequest,
		})
	}

	// Parse month and year or use current
	now := time.Now().UTC()
	month := c.QueryInt("month", int(now.Month()))
	year := c.QueryInt("year", now.Year())

	if month < 1 || month > 12 {
		return c.Status(http.StatusBadRequest).JSON(Response{
			Message: "Month must be between 1 and 12",
			Code:    http.StatusBadRequest,
		})
	}

	monthStart := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	monthEnd := monthStart.AddDate(0, 1, 0)

	// Get weekly breakdown
	var weeklyStats []WeeklyStats
	totalPomodoros := 0
	totalTasks := 0
	totalTime := 0

	current := monthStart
	for current.Before(monthEnd) {
		// Find the start of the week (Monday)
		weekday := int(current.Weekday())
		if weekday == 0 {
			weekday = 7
		}
		weekStart := current.AddDate(0, 0, -weekday+1)
		weekEndDate := weekStart.Add(7 * 24 * time.Hour)

		// Adjust if week extends beyond month
		if weekEndDate.After(monthEnd) {
			weekEndDate = monthEnd
		}

		// Get sessions for this week within the month
		sessionsColl := db.GetDBCollection("sessions")
		sessionsFilter := bson.M{
			"user_id":    userObjectID,
			"status":     model.SessionCompleted,
			"type":       model.Focus,
			"started_at": bson.M{"$gte": weekStart, "$lt": weekEndDate},
		}

		cursor, _ := sessionsColl.Find(c.Context(), sessionsFilter)
		var sessions []model.Session
		cursor.All(c.Context(), &sessions)

		// Get completed tasks for this week within the month
		tasksColl := db.GetDBCollection("tasks")
		tasksFilter := bson.M{
			"user_id":    userObjectID,
			"status":     model.TaskCompleted,
			"updated_at": bson.M{"$gte": weekStart, "$lt": weekEndDate},
		}

		completedTasksCount, _ := tasksColl.CountDocuments(c.Context(), tasksFilter)

		// Calculate weekly total time
		weekTotalTime := 0
		for _, session := range sessions {
			weekTotalTime += int(session.Duration)
		}

		weeklyStats = append(weeklyStats, WeeklyStats{
			WeekStart:          weekStart.Format("2006-01-02"),
			WeekEnd:            weekEndDate.Add(-time.Second).Format("2006-01-02"),
			CompletedPomodoros: len(sessions),
			CompletedTasks:     int(completedTasksCount),
			TotalTime:          weekTotalTime,
		})

		totalPomodoros += len(sessions)
		totalTasks += int(completedTasksCount)
		totalTime += weekTotalTime

		current = weekEndDate
	}

	stats := MonthlyStats{
		Month:              time.Month(month).String(),
		Year:               year,
		CompletedPomodoros: totalPomodoros,
		CompletedTasks:     totalTasks,
		TotalTime:          totalTime,
		WeeklyBreakdown:    weeklyStats,
	}

	return c.Status(http.StatusOK).JSON(Response{
		Message: "Monthly statistics retrieved successfully",
		Code:    http.StatusOK,
		Data:    stats,
	})
}
