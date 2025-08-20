package units

import "testing"

func TestSanitize_CritChance_And_DamageAmp(t *testing.T) {
	t.Parallel()

	// CritChance negative → 0; >1 stays (no gameplay cap here).
	s := build(t, WithCritChance(-0.3))
	if s.Offense.CritChance != 0 {
		t.Fatalf("CritChance negative must clamp to 0, got %v", s.Offense.CritChance)
	}
	s2 := build(t, WithCritChance(1.7))
	if s2.Offense.CritChance != 1.7 {
		t.Fatalf("CritChance should not be capped in data layer, got %v", s2.Offense.CritChance)
	}

	// DamageAmp can be negative (do not clamp).
	s3 := build(t, WithDamageAmp(-0.25))
	if s3.Offense.DamageAmp != -0.25 {
		t.Fatalf("DamageAmp should keep negative values, got %v", s3.Offense.DamageAmp)
	}
}
