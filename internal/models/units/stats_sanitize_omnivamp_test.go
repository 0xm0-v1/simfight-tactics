package units

import (
	"math"
	"testing"
)

func TestSanitize_Omnivamp_MinLEMax_AndClampCurrent(t *testing.T) {
	t.Parallel()

	s := build(t, WithOmnivampValues(0.10, 0.05, 0.20)) // min>max, current>min
	if s.Offense.Omnivamp.OmnivampMin != 0.10 {
		t.Fatalf("min expected 0.10, got %v", s.Offense.Omnivamp.OmnivampMin)
	}
	if s.Offense.Omnivamp.OmnivampMax != 0.10 {
		t.Fatalf("max should be raised to min=0.10, got %v", s.Offense.Omnivamp.OmnivampMax)
	}
	if s.Offense.Omnivamp.CurrentOmnivamp != 0.10 {
		t.Fatalf("current should be clamped to [min,max]=0.10, got %v", s.Offense.Omnivamp.CurrentOmnivamp)
	}
}

func TestValidate_Omnivamp_Current_NonFinite_Fails(t *testing.T) {
	t.Parallel()

	_, err := NewStats(
		WithRange(1),
		WithOmnivampValues(0, 0.2, math.NaN()),
	)
	if err == nil {
		t.Fatalf("expected error on non-finite CurrentOmnivamp")
	}
}
