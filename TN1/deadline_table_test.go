package main

import "testing"

// Test tabulaire pour DaysUntilDeadline.
func TestDaysUntilDeadline_TableDriven(t *testing.T) {
	tests := []struct {
		name        string
		currentDate string
		deadline    string
		wantDays    int
		wantErr     string
	}{
		{
			name:        "positif_echeance_future",
			currentDate: "2025-05-26",
			deadline:    "2025-06-01",
			wantDays:    6,
			wantErr:     "",
		},
		{
			name:        "positif_meme_jour",
			currentDate: "2025-05-26",
			deadline:    "2025-05-26",
			wantDays:    0,
			wantErr:     "",
		},
		{
			name:        "negatif_format_date_actuelle_invalide",
			currentDate: "2025-13-01",
			deadline:    "2025-05-26",
			wantDays:    0,
			wantErr:     "invalid current date format",
		},
		{
			name:        "negatif_deadline_avant_date_actuelle",
			currentDate: "2025-05-26",
			deadline:    "2025-05-25",
			wantDays:    0,
			wantErr:     "deadline cannot be before current date",
		},
		{
			name:        "negatif_date_mal_formatee",
			currentDate: "2025/05/26",
			deadline:    "2025-06-01",
			wantDays:    0,
			wantErr:     "invalid current date format",
		},
		{
			name:        "negatif_format_deadline_invalide",
			currentDate: "2025-05-26",
			deadline:    "2025-13-01",
			wantDays:    0,
			wantErr:     "invalid deadline format",
		},
		{
			name:        "edge_annee_bissextile_valide",
			currentDate: "2024-02-28",
			deadline:    "2024-02-29",
			wantDays:    1,
			wantErr:     "",
		},
		{
			name:        "edge_date_impossible_non_bissextile",
			currentDate: "2025-02-29",
			deadline:    "2025-03-01",
			wantDays:    0,
			wantErr:     "invalid current date format",
		},
		{
			name:        "edge_fin_de_mois",
			currentDate: "2025-01-31",
			deadline:    "2025-02-01",
			wantDays:    1,
			wantErr:     "",
		},
		{
			name:        "edge_changement_annee",
			currentDate: "2025-12-31",
			deadline:    "2026-01-01",
			wantDays:    1,
			wantErr:     "",
		},
		{
			name:        "edge_date_vide",
			currentDate: "",
			deadline:    "2025-06-01",
			wantDays:    0,
			wantErr:     "invalid current date format",
		},
		{
			name:        "edge_date_avec_espaces",
			currentDate: " 2025-05-26 ",
			deadline:    "2025-06-01",
			wantDays:    0,
			wantErr:     "invalid current date format",
		},
		{
			name:        "negatif_deadline_vide",
			currentDate: "2025-05-26",
			deadline:    "",
			wantDays:    0,
			wantErr:     "invalid deadline format",
		},
		{
			name:        "negatif_deadline_avec_espaces",
			currentDate: "2025-05-26",
			deadline:    " 2025-06-01 ",
			wantDays:    0,
			wantErr:     "invalid deadline format",
		},
		{
			name:        "positif_grande_difference",
			currentDate: "2025-01-01",
			deadline:    "2026-01-01",
			wantDays:    365,
			wantErr:     "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gotDays, err := DaysUntilDeadline(tc.currentDate, tc.deadline)

			if tc.wantErr == "" {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
			} else {
				if err == nil {
					t.Errorf("expected error %q, got nil", tc.wantErr)
				} else if err.Error() != tc.wantErr {
					t.Errorf("expected error %q, got %q", tc.wantErr, err.Error())
				}
			}

			if gotDays != tc.wantDays {
				t.Errorf("expected %d days, got %d", tc.wantDays, gotDays)
			}
		})
	}
}
