package main

import "testing"

// ============================================================================
// SECTION 1 : Cas positifs — dates valides, calcul nominal
// ============================================================================

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

// Objectif : Cas positif — même jour (vérifie que le résultat est 0).
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

// Objectif : Cas positif — très grande plage de 10 ans (vérifie la précision du calcul).
func TestDaysUntilDeadline_TenYears(t *testing.T) {
	current, due := "2020-01-01", "2030-01-01"
	want := 3653 // 10 ans incluant 3 années bissextiles (2020, 2024, 2028)
	days, err := DaysUntilDeadline(current, due)
	if err != nil {
		t.Fatalf("DaysUntilDeadline(%q,%q) unexpected error: %v", current, due, err)
	}
	if days != want {
		t.Errorf("DaysUntilDeadline(%q,%q) = %d, want %d", current, due, days, want)
	}
}

// ============================================================================
// SECTION 2 : Cas limites — transitions de dates (fin de mois, année, bissextile)
// ============================================================================

// Objectif : Cas limite — année bissextile (vérifie que 29 février est valide).
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

// Objectif : Cas limite — passage de fin de mois (vérifie transition de mois).
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

// Objectif : Cas limite — passage de fin d'année (vérifie transition d'année).
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

// ============================================================================
// SECTION 3 : Cas négatifs — format de currentDate invalide
// ============================================================================

// Objectif : Cas négatif — mois hors plage (vérifie erreur de parsing).
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

// Objectif : Cas négatif — date impossible en année non bissextile (vérifie erreur).
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

// Objectif : Cas négatif — jour hors plage pour le mois (avril a 30 jours, pas 31).
func TestDaysUntilDeadline_DayOutOfRange(t *testing.T) {
	current, due := "2025-04-31", "2025-05-15"
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

// Objectif : Cas négatif — mois zéro (vérifie que "00" est rejeté par le parser).
func TestDaysUntilDeadline_ZeroMonth(t *testing.T) {
	current, due := "2025-00-15", "2025-06-01"
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

// Objectif : Cas négatif — jour zéro (vérifie que le jour "00" est rejeté).
func TestDaysUntilDeadline_ZeroDay(t *testing.T) {
	current, due := "2025-05-00", "2025-06-01"
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

// Objectif : Cas négatif — date avec heure ISO (vérifie que le surplus est rejeté).
func TestDaysUntilDeadline_DateWithTime(t *testing.T) {
	current, due := "2025-05-26T12:00:00", "2025-06-01"
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

// Objectif : Cas négatif — mois/jour sans zéro initial (vérifie la strictness du format).
func TestDaysUntilDeadline_NoLeadingZero(t *testing.T) {
	current, due := "2025-5-1", "2025-06-01"
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

// Objectif : Cas négatif — entrée complètement garbage (vérifie erreur de parsing).
func TestDaysUntilDeadline_GarbageInput(t *testing.T) {
	current, due := "abc", "not-a-date"
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

// ============================================================================
// SECTION 4 : Cas négatifs — format de deadline invalide
// ============================================================================

// Objectif : Cas négatif — format invalide pour la deadline (mois hors plage).
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

// Objectif : Cas négatif — jour hors plage côté deadline (juin a 30 jours, pas 31).
func TestDaysUntilDeadline_DeadlineDayOutOfRange(t *testing.T) {
	current, due := "2025-06-01", "2025-06-31"
	wantErr := "invalid deadline format"
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

// ============================================================================
// SECTION 5 : Cas négatifs — erreur métier (deadline antérieure)
// ============================================================================

// Objectif : Cas négatif — échéance antérieure à la date actuelle (vérifie erreur métier).
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

// ============================================================================
// SECTION 6 : Cas négatifs — les deux entrées invalides
// ============================================================================

// Objectif : Cas négatif — les deux entrées invalides simultanément (vérifie que la première erreur est retournée).
func TestDaysUntilDeadline_BothInvalid(t *testing.T) {
	current, due := "invalid", "also-invalid"
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
