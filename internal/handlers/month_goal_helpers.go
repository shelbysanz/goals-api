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
