// Ce fichier contient des benchmarks (tests de performance) pour la fonction DaysUntilDeadline

import "testing"

// Scénario : dates valides (deadline après la date courante).
func BenchmarkDaysUntilDeadline_ValidDates(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = DaysUntilDeadline("2025-05-26", "2025-06-01")
	}
}

// Scénario : même jour (date courante = deadline).
func BenchmarkDaysUntilDeadline_SameDay(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = DaysUntilDeadline("2025-05-26", "2025-05-26")
	}
}

// Scénario : deadline dans le passé.
func BenchmarkDaysUntilDeadline_PastDeadline(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = DaysUntilDeadline("2025-05-26", "2025-05-25")
	}
}

// Scénario : format invalide de la date courante (YYYY/MM/DD).
func BenchmarkDaysUntilDeadline_InvalidCurrentDateFormat(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = DaysUntilDeadline("2025/05/26", "2025-06-01")
	}
}

// Scénario : format invalide de la deadline (YYYY/MM/DD).
func BenchmarkDaysUntilDeadline_InvalidDeadlineFormat(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = DaysUntilDeadline("2025-05-26", "2025/06/01")
	}
}

// Scénario : dates vides.
func BenchmarkDaysUntilDeadline_EmptyDates(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = DaysUntilDeadline("", "")
	}
}
	for i := 0; i < b.N; i++ {
		_, _ = DaysUntilDeadline("", "")
	}
}
