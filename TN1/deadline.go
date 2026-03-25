package main

import (
    "errors"
    "time"
)

// DaysUntilDeadline calculates the number of days between currentDate and deadline.
// Dates must be in "YYYY-MM-DD" format. Returns an error if dates are invalid or deadline is before currentDate.
func DaysUntilDeadline(currentDate, deadline string) (int, error) {
    layout := "2006-01-02"
    current, err := time.Parse(layout, currentDate)
    if err != nil {
        return 0, errors.New("invalid current date format")
    }
    due, err := time.Parse(layout, deadline)
    if err != nil {
        return 0, errors.New("invalid deadline format")
    }
    if due.Before(current) {
        return 0, errors.New("deadline cannot be before current date")
    }
    days := int(due.Sub(current).Hours() / 24)
    return days, nil
}
