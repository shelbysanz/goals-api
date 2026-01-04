package validate

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func ErrInvalidMonthYear() (year, month int, err error) {
	return 0, 0, errors.New("month must be in MM-YYYY format")
}

func ParseMonthYear(input string) (year int, month int, err error) {
	if input == "" {
		return ErrInvalidMonthYear()
	}

	parts := strings.Split(input, "-")
	if len(parts) != 2 {
		return ErrInvalidMonthYear()
	}

	month, monErr := strconv.Atoi(parts[0])
	year, yearErr := strconv.Atoi(parts[1])
	if monErr != err || yearErr != err || month < 1 || month > 12 || year < 2026 || year > 2100 {
		return ErrInvalidMonthYear()
	}

	return year, month, nil
}

func FormatMonthYear(year int, month int) string {
	return fmt.Sprintf("%02d-%04d", month, year)
}
