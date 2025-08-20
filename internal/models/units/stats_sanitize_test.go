package units

import (
	"reflect"
	"testing"
)

// ---------- Test constants (avoid magic numbers) ----------
const (
	// General
	rangeMin1 = 1.0
	zero      = 0.0
	negLarge  = -10.0
	negSmall  = -1.0

	// Offense
	asHighNoCap       = 7.0
	asNeg             = -0.3
	critDamageFloor   = 1.0
	critDamageDefault = 1.4
	critDamageLow     = 0.5
	critDamageHigh    = 1.7
	critChanceDefault = 0.25
	ad55              = 55.0
	ad60              = 60.0
	as090             = 0.9

	// Defense
	armor30 = 30.0

	// Resource presets
	manaMin0    = 0.0
	manaMax60   = 60.0
	manaStart30 = 30.0
	manaRegen1  = 1.0
	manaRegen3  = 3.0
	manaRegen5  = 5.0

	// Resource scenarios
	mStartAboveMaxMin   = 100.0
	mStartAboveMaxMax   = 200.0
	mStartAboveMaxStart = 300.0

	mStartBelowMinMin   = 50.0
	mStartBelowMinMax   = 120.0
	mStartBelowMinStart = 25.0

	mMaxBelowMinMin   = 80.0
	mMaxBelowMinMax   = 30.0
	mMaxBelowMinStart = 80.0

	mNegMinMin   = -10.0
	mNegMinMax   = 40.0
	mNegMinStart = 0.0
	mNegMinRegen = -5.0

	// Omnivamp
	omniNeg   = -0.10
	omniStart = 0.10
)

// ---------- Test helper to express Omnivamp override with the new struct ----------
// NOTE: This helper only updates Offense.Omnivamp (via WithOmnivampValues) and
// does NOT overwrite other Offense fields.
func withOmnivamp(min, max, start float64) Option {
	return WithOmnivampValues(min, max, start)
}

// build enforces a valid minimal Range and returns sanitized Stats built from options.
// WithRange(1) is prepended so validation passes unless an explicit WithRange later overrides it.
func build(t *testing.T, opts ...Option) Stats {
	t.Helper()
	opts = append([]Option{WithRange(rangeMin1)}, opts...)
	s, err := NewStats(opts...)
	if err != nil {
		t.Fatalf("build failed: %v", err)
	}
	return s
}

func TestSanitize_NonNegatives(t *testing.T) {
	t.Parallel()

	s := build(t,
		WithAD(negLarge),
		WithAP(negSmall),
		WithHP(negLarge/2),
		WithArmor(negSmall*2),
		WithMR(negSmall*3),
		WithDurability(negSmall*4),
		withOmnivamp(zero, zero, omniNeg), // will clamp start to 0
	)

	if s.Offense.AD != zero {
		t.Errorf("AD: want %v, got %v", zero, s.Offense.AD)
	}
	if s.Offense.AP != zero {
		t.Errorf("AP: want %v, got %v", zero, s.Offense.AP)
	}
	if s.Defense.HP != zero {
		t.Errorf("HP: want %v, got %v", zero, s.Defense.HP)
	}
	if s.Defense.Armor != zero {
		t.Errorf("Armor: want %v, got %v", zero, s.Defense.Armor)
	}
	if s.Defense.MR != zero {
		t.Errorf("MR: want %v, got %v", zero, s.Defense.MR)
	}
	if s.Defense.Durability != zero {
		t.Errorf("Durability: want %v, got %v", zero, s.Defense.Durability)
	}
	wantOmni := Omnivamp{OmnivampMin: zero, OmnivampMax: zero, CurrentOmnivamp: zero} // clamped to 0
	if !reflect.DeepEqual(s.Offense.Omnivamp, wantOmni) {
		t.Errorf("Omnivamp: want %+v, got %+v", wantOmni, s.Offense.Omnivamp)
	}
}

func TestSanitize_CritDamageFloor(t *testing.T) {
	t.Parallel()

	s := build(t, WithCritDamage(critDamageLow))
	if s.Offense.CritDamage != critDamageFloor {
		t.Errorf("CritDamage: want %v floor, got %v", critDamageFloor, s.Offense.CritDamage)
	}

	s2 := build(t, WithCritDamage(critDamageHigh))
	if s2.Offense.CritDamage != critDamageHigh {
		t.Errorf("CritDamage: want unchanged %v, got %v", critDamageHigh, s2.Offense.CritDamage)
	}
}

func TestSanitize_AttackSpeed(t *testing.T) {
	t.Parallel()

	t.Run("NegativeClampedToZero", func(t *testing.T) {
		s := build(t, WithAS(asNeg))
		if s.Offense.AS != zero {
			t.Errorf("AS: want %v when negative, got %v", zero, s.Offense.AS)
		}
	})

	t.Run("NoGameplayCapApplied", func(t *testing.T) {
		s := build(t, WithAS(asHighNoCap))
		if s.Offense.AS != asHighNoCap {
			t.Errorf("AS: should be unchanged (no cap in data layer), got %v", s.Offense.AS)
		}
	})
}

func TestSanitize_ResourceInvariants(t *testing.T) {
	t.Parallel()

	t.Run("StartAboveMax_ClampedToMax", func(t *testing.T) {
		s := build(t, WithMana(mStartAboveMaxMin, mStartAboveMaxMax, mStartAboveMaxStart, manaRegen5, zero)) // start > max
		if s.Resource.ManaMin != mStartAboveMaxMin || s.Resource.ManaMax != mStartAboveMaxMax || s.Resource.ManaStart != mStartAboveMaxMax || s.Resource.ManaRegen != manaRegen5 {
			t.Errorf("Resource clamp to Max failed: %+v", s.Resource)
		}
	})

	t.Run("StartBelowMin_ClampedToMin", func(t *testing.T) {
		s := build(t, WithMana(mStartBelowMinMin, mStartBelowMinMax, mStartBelowMinStart, manaRegen3, zero)) // start < min
		if s.Resource.ManaStart != mStartBelowMinMin {
			t.Errorf("ManaStart clamp to Min failed: got %v", s.Resource.ManaStart)
		}
	})

	t.Run("MaxBelowMin_AdjustedUpToMin", func(t *testing.T) {
		s := build(t, WithMana(mMaxBelowMinMin, mMaxBelowMinMax, mMaxBelowMinStart, manaRegen1, zero)) // max < min
		if s.Resource.ManaMax != mMaxBelowMinMin {
			t.Errorf("ManaMax bumped to Min failed: got %v", s.Resource.ManaMax)
		}
		if s.Resource.ManaStart != mMaxBelowMinMin {
			t.Errorf("ManaStart re-clamped to new Max failed: got %v", s.Resource.ManaStart)
		}
	})

	t.Run("NegativeMinAndRegen_ClampedToZero", func(t *testing.T) {
		s := build(t, WithMana(mNegMinMin, mNegMinMax, mNegMinStart, mNegMinRegen, zero))
		if s.Resource.ManaMin != zero {
			t.Errorf("ManaMin: want %v, got %v", zero, s.Resource.ManaMin)
		}
		if s.Resource.ManaRegen != zero {
			t.Errorf("ManaRegen: want %v, got %v", zero, s.Resource.ManaRegen)
		}
	})
}

func TestSanitize_BulkVsGranular_Offense(t *testing.T) {
	t.Parallel()

	// Bulk path: full WithOffense block
	a := build(t,
		WithOffense(OffenseStats{
			Range:      rangeMin1, // ensure validation passes
			AD:         ad55,
			Omnivamp:   Omnivamp{OmnivampMin: zero, OmnivampMax: zero, CurrentOmnivamp: omniNeg}, // will clamp to 0
			CritChance: critChanceDefault,
			CritDamage: critDamageDefault,
		}),
	)

	// Granular path: independent setters + omnivamp helper
	b := build(t,
		WithAD(ad55),
		withOmnivamp(zero, zero, omniNeg),
	)

	if !reflect.DeepEqual(a.Offense, b.Offense) {
		t.Errorf("Offense mismatch between bulk and granular.\n bulk=%+v\n gran=%+v", a.Offense, b.Offense)
	}
}

func TestSanitize_Idempotent(t *testing.T) {
	t.Parallel()

	s1 := build(t,
		WithAD(ad60),
		WithAS(as090),
		WithArmor(armor30),
		WithMana(manaMin0, manaMax60, manaStart30, manaRegen5, zero),
		WithCritDamage(critDamageDefault),
		withOmnivamp(zero, zero, omniStart),
	)

	// Re-run normalization with no changes (With() without options).
	s2, err := s1.With()
	if err != nil {
		t.Fatalf("unexpected error on With(): %v", err)
	}

	if !reflect.DeepEqual(s1, s2) {
		t.Errorf("sanitize should be idempotent.\n s1=%+v\n s2=%+v", s1, s2)
	}
}
