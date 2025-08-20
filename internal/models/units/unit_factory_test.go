package units

import "testing"

func TestBuildUnit_NoRole_Error(t *testing.T) {
	t.Parallel()
	cfg := RolesLoader{}
	if _, err := BuildUnit("Foo", 1, nil, nil, cfg, WithRange(1)); err == nil {
		t.Fatalf("expected error for missing role")
	}
}

func TestBuildUnit_NoOverride_ReturnsBase(t *testing.T) {
	t.Parallel()
	cfg := RolesLoader{
		RoleTypes:     []string{"Tank"},
		DamageTypes:   []string{"Attack"},
		StatsPerRoles: map[string]any{}, // no overrides
	}
	u, err := BuildUnit("Foo", 1, nil, []string{"Attack Tank"}, cfg, WithRange(1), WithAD(10))
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if u.Stats.Offense.AD != 10 {
		t.Fatalf("expected base stats with AD=10, got %+v", u.Stats.Offense)
	}
}
