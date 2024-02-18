package domain

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

const longDateForm = "Jan 2, 2006 at 3:04pm (MST)"

func FormatDate(due time.Time) string {
	if reflect.ValueOf(due).IsZero() {
		return ""
	}

	s := due.Format(longDateForm)
	value := strings.Split(s, "at")[0]
	return value
}

func ParseDate(date string) (time.Time, error) {
	if len(strings.Trim(date, " \n")) < 1 {
		return time.Time{}, nil
	}

	input := fmt.Sprintf("%s at 10:00pm (EST)", date)
	t, err := time.Parse(longDateForm, input)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse date")
	}
	return t, nil
}
