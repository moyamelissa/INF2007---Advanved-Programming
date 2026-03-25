package main

import "testing"

// Fuzz test : vérifie que la fonction ne panic jamais et respecte des invariants simples.
func FuzzDaysUntilDeadline_Standalone(f *testing.F) {
	seeds := [][2]string{
		{"2025-05-26", "2025-06-01"},
		{"2025-05-26", "2025-05-26"},
		{"2025-05-26", "2025-05-25"},
		{"2025/05/26", "2025-06-01"},
		{"2025-13-01", "2025-06-01"},
		{"", ""},
	}

	for _, s := range seeds {
		f.Add(s[0], s[1])
	}

	f.Fuzz(func(t *testing.T, currentDate string, deadline string) {
		days, err := DaysUntilDeadline(currentDate, deadline)

		if err != nil {
			if days != 0 {
				t.Errorf("when error occurs, expected 0 days, got %d", days)
			}
			return
		}

		if days < 0 {
			t.Errorf("expected non-negative days, got %d", days)
		}
	})
}
