package main

import "testing"

func BenchmarkDaysUntilDeadline_ValidDates(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = DaysUntilDeadline("2025-05-26", "2025-06-01")
	}
}

func BenchmarkDaysUntilDeadline_SameDay(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = DaysUntilDeadline("2025-05-26", "2025-05-26")
	}
}

func BenchmarkDaysUntilDeadline_PastDeadline(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = DaysUntilDeadline("2025-05-26", "2025-05-25")
	}
}

func BenchmarkDaysUntilDeadline_InvalidCurrentDateFormat(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = DaysUntilDeadline("2025/05/26", "2025-06-01")
	}
}

func BenchmarkDaysUntilDeadline_InvalidDeadlineFormat(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = DaysUntilDeadline("2025-05-26", "2025/06/01")
	}
}

func BenchmarkDaysUntilDeadline_EmptyDates(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = DaysUntilDeadline("", "")
	}
}
