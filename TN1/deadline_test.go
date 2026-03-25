package main

import "testing"

// Test positif : échéance future valide
func TestDaysUntilDeadline_FutureDate(t *testing.T) {
	days, err := DaysUntilDeadline("2025-05-26", "2025-06-01")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if days != 6 {
		t.Errorf("expected 6 days, got %d", days)
	}
}

// Test positif : même jour
func TestDaysUntilDeadline_SameDay(t *testing.T) {
	days, err := DaysUntilDeadline("2025-05-26", "2025-05-26")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if days != 0 {
		t.Errorf("expected 0 days, got %d", days)
	}
}

// Test négatif : format invalide
func TestDaysUntilDeadline_InvalidFormat(t *testing.T) {
	days, err := DaysUntilDeadline("2025-13-01", "2025-05-26")
	if err == nil {
		t.Errorf("expected an error for invalid date format")
	}
	if days != 0 {
		t.Errorf("expected 0 days when error occurs, got %d", days)
	}
	if err != nil && err.Error() != "invalid current date format" {
		t.Errorf("expected error \"invalid current date format\", got %q", err.Error())
	}
}

// Test négatif : deadline avant la date actuelle
func TestDaysUntilDeadline_DeadlineBeforeCurrent(t *testing.T) {
	days, err := DaysUntilDeadline("2025-05-26", "2025-05-25")
	if err == nil {
		t.Errorf("expected an error because deadline is before current date")
	}
	if days != 0 {
		t.Errorf("expected 0 days when error occurs, got %d", days)
	}
	if err != nil && err.Error() != "deadline cannot be before current date" {
		t.Errorf("expected error \"deadline cannot be before current date\", got %q", err.Error())
	}
}

// Test négatif : mauvais séparateur
func TestDaysUntilDeadline_WrongSeparator(t *testing.T) {
	days, err := DaysUntilDeadline("2025/05/26", "2025-06-01")
	if err == nil {
		t.Errorf("expected an error for wrong date separator")
	}
	if days != 0 {
		t.Errorf("expected 0 days when error occurs, got %d", days)
	}
	if err != nil && err.Error() != "invalid current date format" {
		t.Errorf("expected error \"invalid current date format\", got %q", err.Error())
	}
}

// Test négatif : format invalide pour la deadline
func TestDaysUntilDeadline_InvalidDeadlineFormat(t *testing.T) {
	days, err := DaysUntilDeadline("2025-05-26", "2025-13-01")
	if err == nil {
		t.Errorf("expected an error for invalid deadline format")
	}
	if days != 0 {
		t.Errorf("expected 0 days when error occurs, got %d", days)
	}
	if err != nil && err.Error() != "invalid deadline format" {
		t.Errorf("expected error \"invalid deadline format\", got %q", err.Error())
	}
}

// Test limite : année bissextile valide
func TestDaysUntilDeadline_LeapYearValid(t *testing.T) {
	days, err := DaysUntilDeadline("2024-02-28", "2024-02-29")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if days != 1 {
		t.Errorf("expected 1 day, got %d", days)
	}
}

// Test limite : date impossible en année non bissextile
func TestDaysUntilDeadline_InvalidNonLeapDate(t *testing.T) {
	days, err := DaysUntilDeadline("2025-02-29", "2025-03-01")
	if err == nil {
		t.Errorf("expected an error for invalid current date")
	}
	if days != 0 {
		t.Errorf("expected 0 days when error occurs, got %d", days)
	}
	if err != nil && err.Error() != "invalid current date format" {
		t.Errorf("expected error \"invalid current date format\", got %q", err.Error())
	}
}

// Test limite : passage de fin de mois
func TestDaysUntilDeadline_EndOfMonth(t *testing.T) {
	days, err := DaysUntilDeadline("2025-01-31", "2025-02-01")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if days != 1 {
		t.Errorf("expected 1 day, got %d", days)
	}
}

// Test limite : passage de fin d'année
func TestDaysUntilDeadline_EndOfYear(t *testing.T) {
	days, err := DaysUntilDeadline("2025-12-31", "2026-01-01")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if days != 1 {
		t.Errorf("expected 1 day, got %d", days)
	}
}

// Test limite : date actuelle vide
func TestDaysUntilDeadline_EmptyCurrentDate(t *testing.T) {
	days, err := DaysUntilDeadline("", "2025-06-01")
	if err == nil {
		t.Errorf("expected an error for empty current date")
	}
	if days != 0 {
		t.Errorf("expected 0 days when error occurs, got %d", days)
	}
	if err != nil && err.Error() != "invalid current date format" {
		t.Errorf("expected error \"invalid current date format\", got %q", err.Error())
	}
}

// Test limite : date actuelle avec espaces
func TestDaysUntilDeadline_CurrentDateWithSpaces(t *testing.T) {
	days, err := DaysUntilDeadline(" 2025-05-26 ", "2025-06-01")
	if err == nil {
		t.Errorf("expected an error for spaced current date")
	}
	if days != 0 {
		t.Errorf("expected 0 days when error occurs, got %d", days)
	}
	if err != nil && err.Error() != "invalid current date format" {
		t.Errorf("expected error \"invalid current date format\", got %q", err.Error())
	}
}
