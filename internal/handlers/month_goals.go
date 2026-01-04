package handlers

import (
	"net/http"
	"time"

	"goals-api/internal/models"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type MonthGoalHandler struct {
	DB *gorm.DB
}

type CreateMonthGoalRequest struct {
	Month     string `json:"month"`
	Title     string `json:"title"`
	Notes     string `json:"notes"`
	Completed bool   `json:"completed"`
}

type MonthGoalResponse struct {
	ID        uint   `json:"id"`
	Month     string `json:"month"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
	Notes     string `json:"notes"`
}

var ALLOWED_MONTH_GOAL_SORTED_COLUMNS = map[string]string{
	"month":     "month",
	"year":      "year",
	"title":     "title",
	"completed": "completed",
	"createdAt": "created_at",
}

var ALLOWED_SORTED_DIRECTIONS = []string{"asc", "desc"}

func toMonthGoalResponse(m models.MonthGoal) MonthGoalResponse {
	return MonthGoalResponse{
		ID:        m.ID,
		Month:     validate.FormatMonthYear(m.Year, m.Month),
		Title:     m.Title,
		Completed: m.Completed,
		Notes:     m.Notes,
	}
}

func NewMonthGoalHandler(db *gorm.DB) *MonthGoalHandler {
	return &MonthGoalHandler{DB: db}
}

// essentially runs a: SELECT * FROM month_goals;
// marshals results to json and returns an HTTP 200 response
func (mg *MonthGoalHandler) List(c echo.Context) error {
	var goals []models.MonthGoal

	// initialize query
	query := mg.DB

	/* parse and apply query params */
	// optional month filter
	monthParam := c.QueryParam("month")
	if monthParam != "" {
		year, month, err := validate.ParseMonthYear(monthParam)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		query = query.Where("year = ? AND month = ?", year, month)
	}

	if err := query.Find(&goals).Error; err != nil {
		return c.String(http.StatusInternalServerError, "database error")
	}

	response := make([]MonthGoalResponse, 0, len(goals))
	for _, g := range goals {
		response = append(response, toMonthGoalResponse(g))
	}

	return c.JSON(http.StatusOK, goals)
}

func (mg *MonthGoalHandler) Create(c echo.Context) error {
	var req CreateMonthGoalRequest

	if err := c.Bind(&req); err != nil {
		return c.String(http.StatusBadRequest, "invalid JSON")
	}

	year, month, err := validate.ParseMonthYear(req.Month)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	if req.Title == "" {
		return c.String(http.StatusBadRequest, "title is required")
	}

	goal := models.MonthGoal{
		Month:     month,
		Year:      year,
		Title:     req.Title,
		Notes:     req.Notes,
		Completed: req.Completed,
	}

	if err := mg.DB.Create(&goal).Error; err != nil {
		return c.String(http.StatusInternalServerError, "database error")
	}

	return c.JSON(http.StatusCreated, toMonthGoalResponse(goal))
}
