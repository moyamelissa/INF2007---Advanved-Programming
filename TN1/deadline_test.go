package main

import "testing"

// Objectif : Cas positif — échéance future valide (vérifie le calcul nominal).
func TestDaysUntilDeadline_FutureDate(t *testing.T) {
	current, due := "2025-05-26", "2025-06-01"
	want := 6
	days, err := DaysUntilDeadline(current, due)
	if err != nil {
		t.Fatalf("DaysUntilDeadline(%q,%q) unexpected error: %v", current, due, err)
	}
	if days != want {
		t.Errorf("DaysUntilDeadline(%q,%q) = %d, want %d", current, due, days, want)
	}
}

// Objectif : Cas limite positif — même jour (vérifie que le résultat est 0).
func TestDaysUntilDeadline_SameDay(t *testing.T) {
	current, due := "2025-05-26", "2025-05-26"
	want := 0
	days, err := DaysUntilDeadline(current, due)
	if err != nil {
		t.Fatalf("DaysUntilDeadline(%q,%q) unexpected error: %v", current, due, err)
	}
	if days != want {
		t.Errorf("DaysUntilDeadline(%q,%q) = %d, want %d", current, due, days, want)
	}
}

// Objectif : Cas négatif — format de date invalide (vérifie erreur et days==0).
func TestDaysUntilDeadline_InvalidFormat(t *testing.T) {
	current, due := "2025-13-01", "2025-05-26"
	wantErr := "invalid current date format"
	days, err := DaysUntilDeadline(current, due)
	if err == nil {
		t.Fatalf("DaysUntilDeadline(%q,%q) expected error, got nil", current, due)
	}
	if days != 0 {
		t.Errorf("DaysUntilDeadline(%q,%q) = %d, want 0 on error", current, due, days)
	}
	if err.Error() != wantErr {
		t.Errorf("DaysUntilDeadline(%q,%q) error = %q, want %q", current, due, err.Error(), wantErr)
	}
}

// Objectif : Cas négatif — échéance antérieure (vérifie erreur métier et days==0).
func TestDaysUntilDeadline_DeadlineBeforeCurrent(t *testing.T) {
	current, due := "2025-05-26", "2025-05-25"
	wantErr := "deadline cannot be before current date"
	days, err := DaysUntilDeadline(current, due)
	if err == nil {
		t.Fatalf("DaysUntilDeadline(%q,%q) expected error, got nil", current, due)
	}
	if days != 0 {
		t.Errorf("DaysUntilDeadline(%q,%q) = %d, want 0 on error", current, due, days)
	}
	if err.Error() != wantErr {
		t.Errorf("DaysUntilDeadline(%q,%q) error = %q, want %q", current, due, err.Error(), wantErr)
	}
}

// Objectif : Cas négatif — séparateur invalide (vérifie erreur de parsing).
func TestDaysUntilDeadline_WrongSeparator(t *testing.T) {
	current, due := "2025/05/26", "2025-06-01"
	wantErr := "invalid current date format"
	days, err := DaysUntilDeadline(current, due)
	if err == nil {
		t.Fatalf("DaysUntilDeadline(%q,%q) expected error, got nil", current, due)
	}
	if days != 0 {
		t.Errorf("DaysUntilDeadline(%q,%q) = %d, want 0 on error", current, due, days)
	}
	if err.Error() != wantErr {
		t.Errorf("DaysUntilDeadline(%q,%q) error = %q, want %q", current, due, err.Error(), wantErr)
	}
}

// Objectif : Cas négatif — format invalide pour la deadline (vérifie erreur et days==0).
func TestDaysUntilDeadline_InvalidDeadlineFormat(t *testing.T) {
	current, due := "2025-05-26", "2025-13-01"
	wantErr := "invalid deadline format"
	days, err := DaysUntilDeadline(current, due)
	if err == nil {
		t.Fatalf("DaysUntilDeadline(%q,%q) expected error, got nil", current, due)
	}
	if days != 0 {
		t.Fatalf("DaysUntilDeadline(%q,%q) = %d, want 0 on error", current, due, days)
	}
	if err.Error() != wantErr {
		t.Fatalf("DaysUntilDeadline(%q,%q) error = %q, want %q", current, due, err.Error(), wantErr)
	}
}

// Objectif : Cas limite positif — année bissextile (vérifie que 29 février est valide).
func TestDaysUntilDeadline_LeapYearValid(t *testing.T) {
	current, due := "2024-02-28", "2024-02-29"
	want := 1
	days, err := DaysUntilDeadline(current, due)
	if err != nil {
		t.Fatalf("DaysUntilDeadline(%q,%q) unexpected error: %v", current, due, err)
	}
	if days != want {
		t.Errorf("DaysUntilDeadline(%q,%q) = %d, want %d", current, due, days, want)
	}
}

// Objectif : Cas limite négatif — date impossible en année non bissextile (vérifie erreur).
func TestDaysUntilDeadline_InvalidNonLeapDate(t *testing.T) {
	current, due := "2025-02-29", "2025-03-01"
	wantErr := "invalid current date format"
	days, err := DaysUntilDeadline(current, due)
	if err == nil {
		t.Fatalf("DaysUntilDeadline(%q,%q) expected error, got nil", current, due)
	}
	if days != 0 {
		t.Errorf("DaysUntilDeadline(%q,%q) = %d, want 0 on error", current, due, days)
	}
	if err.Error() != wantErr {
		t.Errorf("DaysUntilDeadline(%q,%q) error = %q, want %q", current, due, err.Error(), wantErr)
	}
}

// Objectif : Cas limite positif — passage de fin de mois (vérifie transition de mois).
func TestDaysUntilDeadline_EndOfMonth(t *testing.T) {
	current, due := "2025-01-31", "2025-02-01"
	want := 1
	days, err := DaysUntilDeadline(current, due)
	if err != nil {
		t.Fatalf("DaysUntilDeadline(%q,%q) unexpected error: %v", current, due, err)
	}
	if days != want {
		t.Fatalf("DaysUntilDeadline(%q,%q) = %d, want %d", current, due, days, want)
	}
}

// Objectif : Cas limite positif — passage de fin d'année (vérifie transition d'année).
func TestDaysUntilDeadline_EndOfYear(t *testing.T) {
	current, due := "2025-12-31", "2026-01-01"
	want := 1

	days, err := DaysUntilDeadline(current, due)
	if err != nil {
		t.Fatalf("DaysUntilDeadline(%q,%q) unexpected error: %v", current, due, err)
	}
	if days != want {
		t.Fatalf("DaysUntilDeadline(%q,%q) = %d, want %d", current, due, days, want)
	}
}

// Objectif : Cas négatif — date actuelle vide (vérifie erreur et days==0).
func TestDaysUntilDeadline_EmptyCurrentDate(t *testing.T) {
	current, due := "", "2025-06-01"
	wantErr := "invalid current date format"
	days, err := DaysUntilDeadline(current, due)
	if err == nil {
		t.Fatalf("DaysUntilDeadline(%q,%q) expected error, got nil", current, due)
	}
	if days != 0 {
		t.Fatalf("DaysUntilDeadline(%q,%q) = %d, want 0 on error", current, due, days)
	}
	if err.Error() != wantErr {
		t.Fatalf("DaysUntilDeadline(%q,%q) error = %q, want %q", current, due, err.Error(), wantErr)
	}
}

// Objectif : Cas négatif — date actuelle avec espaces (vérifie erreur de format).
func TestDaysUntilDeadline_CurrentDateWithSpaces(t *testing.T) {
	current, due := " 2025-05-26 ", "2025-06-01"
	wantErr := "invalid current date format"
	days, err := DaysUntilDeadline(current, due)
	if err == nil {
		t.Fatalf("DaysUntilDeadline(%q,%q) expected error, got nil", current, due)
	}
	if days != 0 {
		t.Fatalf("DaysUntilDeadline(%q,%q) = %d, want 0 on error", current, due, days)
	}
	if err.Error() != wantErr {
		t.Fatalf("DaysUntilDeadline(%q,%q) error = %q, want %q", current, due, err.Error(), wantErr)
	}
}

// Objectif : Cas négatif — deadline vide (vérifie erreur et days==0).
func TestDaysUntilDeadline_EmptyDeadline(t *testing.T) {
	current, due := "2025-05-26", ""
	wantErr := "invalid deadline format"
	days, err := DaysUntilDeadline(current, due)
	if err == nil {
		t.Fatalf("DaysUntilDeadline(%q,%q) expected error, got nil", current, due)
	}
	if days != 0 {
		t.Fatalf("DaysUntilDeadline(%q,%q) = %d, want 0 on error", current, due, days)
	}
	if err.Error() != wantErr {
		t.Fatalf("DaysUntilDeadline(%q,%q) error = %q, want %q", current, due, err.Error(), wantErr)
	}
}

// Objectif : Cas négatif — deadline avec espaces (vérifie erreur de format).
func TestDaysUntilDeadline_DeadlineWithSpaces(t *testing.T) {
	current, due := "2025-05-26", " 2025-06-01 "
	wantErr := "invalid deadline format"
	days, err := DaysUntilDeadline(current, due)
	if err == nil {
		t.Fatalf("DaysUntilDeadline(%q,%q) expected error, got nil", current, due)
	}
	if days != 0 {
		t.Fatalf("DaysUntilDeadline(%q,%q) = %d, want 0 on error", current, due, days)
	}
	if err.Error() != wantErr {
		t.Fatalf("DaysUntilDeadline(%q,%q) error = %q, want %q", current, due, err.Error(), wantErr)
	}
}

// Objectif : Cas positif — grande différence de dates (vérifie calcul sur 365 jours).
func TestDaysUntilDeadline_LargeDateRange(t *testing.T) {
	current, due := "2025-01-01", "2026-01-01"
	want := 365
	days, err := DaysUntilDeadline(current, due)
	if err != nil {
		t.Fatalf("DaysUntilDeadline(%q,%q) unexpected error: %v", current, due, err)
	}
	if days != want {
		t.Errorf("DaysUntilDeadline(%q,%q) = %d, want %d", current, due, days, want)
	}
}
