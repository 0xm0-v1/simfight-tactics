package units

import (
	"reflect"
	"testing"
)

// build is a test helper that always enforces the champion Range (>=1) so validation passes,
// then returns sanitized Stats built from the provided options.
func build(t *testing.T, opts ...Option) Stats {
	t.Helper()
	// Put WithRange(1) FIRST so an explicit WithRange(...) in opts can override it.
	opts = append([]Option{WithRange(1)}, opts...)
	s, err := NewStats(opts...)
	if err != nil {
		t.Fatalf("build failed: %v", err)
	}
	return s
}

func TestSanitize_NonNegatives(t *testing.T) {
	t.Parallel()

	s := build(t,
		WithAD(-10),
		WithAP(-1),
		WithHP(-5),
		WithArmor(-2),
		WithMR(-3),
		WithDurability(-4),
		WithOmnivamp(-0.25),
	)

	if s.Offense.AD != 0 {
		t.Errorf("AD: want 0, got %v", s.Offense.AD)
	}
	if s.Offense.AP != 0 {
		t.Errorf("AP: want 0, got %v", s.Offense.AP)
	}
	if s.Defense.HP != 0 {
		t.Errorf("HP: want 0, got %v", s.Defense.HP)
	}
	if s.Defense.Armor != 0 {
		t.Errorf("Armor: want 0, got %v", s.Defense.Armor)
	}
	if s.Defense.MR != 0 {
		t.Errorf("MR: want 0, got %v", s.Defense.MR)
	}
	if s.Defense.Durability != 0 {
		t.Errorf("Durability: want 0, got %v", s.Defense.Durability)
	}
	if s.Offense.Omnivamp != 0 {
		t.Errorf("Omnivamp: want 0, got %v", s.Offense.Omnivamp)
	}
}

func TestSanitize_CritDamageFloor(t *testing.T) {
	t.Parallel()

	s := build(t, WithCritDamage(0.5))
	if s.Offense.CritDamage != 1.0 {
		t.Errorf("CritDamage: want 1.0 floor, got %v", s.Offense.CritDamage)
	}

	s2 := build(t, WithCritDamage(1.7))
	if s2.Offense.CritDamage != 1.7 {
		t.Errorf("CritDamage: want unchanged 1.7, got %v", s2.Offense.CritDamage)
	}
}

func TestSanitize_AttackSpeed(t *testing.T) {
	t.Parallel()

	t.Run("NegativeClampedToZero", func(t *testing.T) {
		s := build(t, WithAS(-0.3))
		if s.Offense.AS != 0 {
			t.Errorf("AS: want 0 when negative, got %v", s.Offense.AS)
		}
	})

	t.Run("NoGameplayCapApplied", func(t *testing.T) {
		s := build(t, WithAS(7.0))
		if s.Offense.AS != 7.0 {
			t.Errorf("AS: should be unchanged (no cap in data layer), got %v", s.Offense.AS)
		}
	})
}

func TestSanitize_ResourceInvariants(t *testing.T) {
	t.Parallel()

	t.Run("StartAboveMax_ClampedToMax", func(t *testing.T) {
		s := build(t, WithMana(100, 200, 300, 5)) // start > max
		if s.Resource.ManaMin != 100 || s.Resource.ManaMax != 200 || s.Resource.ManaStart != 200 || s.Resource.ManaRegen != 5 {
			t.Errorf("Resource clamp to Max failed: %+v", s.Resource)
		}
	})

	t.Run("StartBelowMin_ClampedToMin", func(t *testing.T) {
		s := build(t, WithMana(50, 120, 25, 3)) // start < min
		if s.Resource.ManaStart != 50 {
			t.Errorf("ManaStart clamp to Min failed: got %v", s.Resource.ManaStart)
		}
	})

	t.Run("MaxBelowMin_AdjustedUpToMin", func(t *testing.T) {
		s := build(t, WithMana(80, 30, 80, 1)) // max < min
		if s.Resource.ManaMax != 80 {
			t.Errorf("ManaMax bumped to Min failed: got %v", s.Resource.ManaMax)
		}
		if s.Resource.ManaStart != 80 {
			t.Errorf("ManaStart re-clamped to new Max failed: got %v", s.Resource.ManaStart)
		}
	})

	t.Run("NegativeMinAndRegen_ClampedToZero", func(t *testing.T) {
		s := build(t, WithMana(-10, 40, 0, -5))
		if s.Resource.ManaMin != 0 {
			t.Errorf("ManaMin: want 0, got %v", s.Resource.ManaMin)
		}
		if s.Resource.ManaRegen != 0 {
			t.Errorf("ManaRegen: want 0, got %v", s.Resource.ManaRegen)
		}
	})
}

func TestSanitize_BulkVsGranular_Offense(t *testing.T) {
	t.Parallel()

	a := build(t,
		WithOffense(OffenseStats{
			Range:      1, // pour passer Validate()
			AD:         55,
			Omnivamp:   -0.10,                        // sera clampé à 0
			CritChance: Default().Offense.CritChance, // 0.25
			CritDamage: Default().Offense.CritDamage, // 1.4
		}),
	)
	b := build(t,
		WithAD(55),
		WithOmnivamp(-0.10),
	)

	if !reflect.DeepEqual(a.Offense, b.Offense) {
		t.Errorf("Offense mismatch between bulk and granular.\n bulk=%+v\n gran=%+v", a.Offense, b.Offense)
	}
}

func TestSanitize_Idempotent(t *testing.T) {
	t.Parallel()

	s1 := build(t,
		WithAD(60),
		WithAS(0.9),
		WithArmor(30),
		WithMana(0, 60, 30, 5),
		WithCritDamage(1.4),
		WithOmnivamp(0.1),
	)

	// Re-run normalization without changes (With() no options).
	s2, err := s1.With() // applies Validate() + normalized() again
	if err != nil {
		t.Fatalf("unexpected error on With(): %v", err)
	}

	if !reflect.DeepEqual(s1, s2) {
		t.Errorf("sanitize should be idempotent.\n s1=%+v\n s2=%+v", s1, s2)
	}
}
