package utils

import (
	"errors"
	"strconv"
	"time"

	"gopkg.in/guregu/null.v3"
)

// FloatToString converts float64 to string
func FloatToString(input float64) string {
	return strconv.FormatFloat(input, 'f', -1, 64)
}

// StringToFloat converts string to float64
func StringToFloat(input string) float64 {
	value, err := strconv.ParseFloat(input, 64)
	if err != nil {
		return 0
	}

	return value
}

// NullStringToString converts null.String to string, with empty string as default value if it is not valid
func NullStringToString(value null.String) string {
	if value.Valid {
		return value.String
	}

	return ""
}

// NullIntToInt64 converts null.Int to int64, with 0 as default value if it is not valid
func NullIntToInt64(value null.Int) int64 {
	if value.Valid {
		return value.Int64
	}

	return int64(0)
}

// NullTimeToTime converts null.Int to Time
func NullTimeToTime(value null.Time) time.Time {
	if value.Valid {
		return value.Time
	}

	return time.Time{}
}

// ParseTimeFromString parse string into Time
func ParseTimeFromString(value string) (time.Time, error) {
	res, err := time.Parse("2006-01-02 15:04:05", value)
	if err != nil {
		return time.Time{}, errors.New("Parse Time Failed")
	}

	return res, nil
}
