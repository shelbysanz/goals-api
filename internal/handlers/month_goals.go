package handlers

import (
	"net/http"

	"goals-api/internal/models"
	"goals-api/internal/validate"

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

type UpdateMonthGoalRequest struct {
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
	CreatedAt string `json:"created_at"`
}

func NewMonthGoalHandler(db *gorm.DB) *MonthGoalHandler {
	return &MonthGoalHandler{DB: db}
}

// essentially runs: SELECT * FROM month_goals;
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

	// optional sort
	sortParam := c.QueryParam("sort")
	if orderBy, ok := parseSortParam(sortParam); ok {
		query = query.Order(orderBy)
	}

	if err := query.Find(&goals).Error; err != nil {
		return c.String(http.StatusInternalServerError, "database error")
	}

	response := make([]MonthGoalResponse, 0, len(goals))
	for _, g := range goals {
		response = append(response, toMonthGoalResponse(g))
	}

	return c.JSON(http.StatusOK, response)
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
