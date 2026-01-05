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
	Month     *string `json:"month"`
	Title     *string `json:"title"`
	Notes     *string `json:"notes"`
	Completed *bool   `json:"completed"`
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

func (mg *MonthGoalHandler) Update(c echo.Context) error {
	id := c.Param("id") // note this is param not query param
	if id == "" {
		return c.String(http.StatusBadRequest, "id is required")
	}

	var req UpdateMonthGoalRequest
	if err := c.Bind(&req); err != nil {
		return c.String(http.StatusBadRequest, "invalid JSON")
	}

	var goal models.MonthGoal
	if err := mg.DB.First(&goal, id).Error; err != nil {
		return c.String(http.StatusNotFound, "month goal not found")
	}

	patch, err := buildMonthGoalPatch(req)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	applyMonthGoalPatch(&goal, patch)

	if err := mg.DB.Save(&goal).Error; err != nil {
		return c.String(http.StatusInternalServerError, "database error")
	}

	return c.JSON(http.StatusOK, toMonthGoalResponse(goal))
}

func (mg *MonthGoalHandler) Delete(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.String(http.StatusBadRequest, "id is required")
	}

	result := mg.DB.Delete(&models.MonthGoal{}, id)
	if result.Error != nil {
		return c.String(http.StatusInternalServerError, "database error")
	}

	if result.RowsAffected == 0 {
		return c.String(http.StatusNotFound, "month goal not found")
	}

	return c.NoContent(http.StatusNoContent)
}
