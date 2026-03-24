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
	_, err := DaysUntilDeadline("2025-13-01", "2025-05-26")
	if err == nil {
		t.Errorf("expected an error for invalid date format")
	}
}

// Test négatif : deadline avant la date actuelle
func TestDaysUntilDeadline_DeadlineBeforeCurrent(t *testing.T) {
	_, err := DaysUntilDeadline("2025-05-26", "2025-05-25")
	if err == nil {
		t.Errorf("expected an error because deadline is before current date")
	}
}

// Test négatif : mauvais séparateur
func TestDaysUntilDeadline_WrongSeparator(t *testing.T) {
	_, err := DaysUntilDeadline("2025/05/26", "2025-06-01")
	if err == nil {
		t.Errorf("expected an error for wrong date separator")
	}
}

// Test négatif : format invalide pour la deadline
func TestDaysUntilDeadline_InvalidDeadlineFormat(t *testing.T) {
	_, err := DaysUntilDeadline("2025-05-26", "2025-13-01")
	if err == nil {
		t.Errorf("expected an error for invalid deadline format")
	}
}
