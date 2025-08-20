package units

import "testing"

func TestApplyRoleMapToStats_ReportsUnknownAndTypeErrors(t *testing.T) {
	t.Parallel()

	roleDoc := map[string]any{
		"offense": map[string]any{
			"attack_speed": "fast", // type err
			"unknown_key":  123.0,  // unknown
		},
		"defense": 123, // type err: object expected
	}
	var dst Stats
	report := applyRoleMapToStats(&dst, roleDoc)
	if len(report.UnknownKeys) == 0 || len(report.TypeErrors) == 0 {
		t.Fatalf("expected both unknown keys and type errors, got %+v", report)
	}
}

func TestValidateRolesConfig_Basic(t *testing.T) {
	t.Parallel()

	cfg := RolesLoader{
		RoleTypes:   []string{"Tank", "Fighter"},
		DamageTypes: []string{"Attack", "Magic"},
		StatsPerRoles: map[string]any{
			"TankX": map[string]any{}, // unknown role
			"Tank":  42,               // non-object
			"Fighter": map[string]any{
				"offense": map[string]any{
					"unknown":      1.0,    // unknown key
					"attack_speed": "oops", // type error
				},
			},
		},
	}
	if err := ValidateRolesConfig(cfg); err == nil {
		t.Fatalf("expected validation error")
	}
}
