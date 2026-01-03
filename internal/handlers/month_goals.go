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

func NewMonthGoalHandler(db *gorm.DB) *MonthGoalHandler {
	return &MonthGoalHandler{DB: db}
}

// essentially runs a: SELECT * FROM month_goals;
// marshals results to json and returns an HTTP 200 response
func (mg *MonthGoalHandler) List(c echo.Context) error {
	var goals []models.MonthGoal

	if err := mg.DB.Find(&goals).Error; err != nil {
		return c.String(http.StatusInternalServerError, "database error")
	}

	return c.JSON(http.StatusOK, goals)
}

func (mg *MonthGoalHandler) Create(c echo.Context) error {
	var req CreateMonthGoalRequest

	if err := c.Bind(&req); err != nil {
		return c.String(http.StatusBadRequest, "invalid JSON")
	}

	if req.Month == "" || req.Title == "" {
		return c.String(http.StatusBadRequest, "both month and title are required")
	}

	month, err := time.Parse("01-2006", req.Month)
	if err != nil {
		return c.String(http.StatusBadRequest, "month must be in MM-YYYY format")
	}

	goal := models.MonthGoal{
		Month:     month,
		Title:     req.Title,
		Notes:     req.Notes,
		Completed: req.Completed,
	}

	if err := mg.DB.Create(&goal).Error; err != nil {
		return c.String(http.StatusInternalServerError, "database error")
	}

	return c.JSON(http.StatusCreated, goal)
}
