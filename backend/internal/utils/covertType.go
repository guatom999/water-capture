package utils

import (
	"strconv"
	"time"
)

func ConvertStringToFloat64(value string) float64 {
	afterConv, _ := strconv.ParseFloat(value, 64)
	return afterConv
}

func ConvertStringToTime(value string) time.Time {
	loc, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		return time.Time{}
	}

	layout := "2006-01-02 15:04"
	parsedTime, err := time.ParseInLocation(layout, value, loc)
	if err != nil {
		return time.Time{}
	}

	return parsedTime
}
