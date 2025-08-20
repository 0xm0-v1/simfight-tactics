package units

import "testing"

func TestStatsForRole_StrictMode_ErrorsOnReport(t *testing.T) {
	t.Parallel()

	cfg := RolesLoader{
		Strict:      true,
		RoleTypes:   []string{"Tank"},
		DamageTypes: []string{"Attack"},
		StatsPerRoles: map[string]any{
			"tank": map[string]any{
				"offense": map[string]any{
					"unknown": 1.0, // should trigger report → error in Strict mode
				},
			},
		},
	}
	_, err := StatsForRole("Attack Tank", cfg, WithRange(1))
	if err == nil {
		t.Fatalf("expected error in strict mode")
	}
}
