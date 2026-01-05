package handlers

import (
	"slices"
	"strings"

	"goals-api/internal/models"
	"goals-api/internal/validate"
)

var ALLOWED_MONTH_GOAL_SORTED_COLUMNS = map[string]string{
	"month":      "month",
	"year":       "year",
	"title":      "title",
	"completed":  "completed",
	"created_at": "created_at",
}

var ALLOWED_SORTED_DIRECTIONS = []string{"asc", "desc"}

type MonthGoalPatch struct {
	Completed *bool
	Year      *int
	Month     *int
	Notes     *string
	Title     *string
}

func toMonthGoalResponse(m models.MonthGoal) MonthGoalResponse {
	return MonthGoalResponse{
		ID:        m.ID,
		Month:     validate.FormatMonthYear(m.Year, m.Month),
		Title:     m.Title,
		Completed: m.Completed,
		Notes:     m.Notes,
		CreatedAt: m.CreatedAt.Format("01-02-2006 at 03:04 pm"),
	}
}

func parseSortParam(sort string) (string, bool) {
	if sort == "" {
		return "", false
	}

	parts := strings.Split(sort, ":")
	if len(parts) != 2 {
		return "", false
	}

	col := parts[0]
	dir := parts[1]

	column, ok := ALLOWED_MONTH_GOAL_SORTED_COLUMNS[col]
	if !ok {
		return "", false
	}

	if !slices.Contains(ALLOWED_SORTED_DIRECTIONS, dir) {
		return "", false
	}

	return column + " " + dir, true
}

func buildMonthGoalPatch(req UpdateMonthGoalRequest) (*MonthGoalPatch, error) {
	patch := &MonthGoalPatch{}

	if req.Completed != nil {
		patch.Completed = req.Completed
	}

	if req.Month != nil {
		month, year, err := validate.ParseMonthYear(*req.Month)
		if err != nil {
			return nil, err
		}
		patch.Month = &month
		patch.Year = &year
	}

	if req.Notes != nil {
		patch.Notes = req.Notes
	}

	if req.Title != nil {
		patch.Title = req.Title
	}

	return patch, nil
}

func applyMonthGoalPatch(goal *models.MonthGoal, patch *MonthGoalPatch) {
	if patch.Completed != nil {
		goal.Completed = *patch.Completed
	}
	if patch.Month != nil && patch.Year != nil {
		goal.Month = *patch.Month
		goal.Year = *patch.Year
	}
	if patch.Notes != nil {
		goal.Notes = *patch.Notes
	}
	if patch.Title != nil {
		goal.Title = *patch.Title
	}
}
