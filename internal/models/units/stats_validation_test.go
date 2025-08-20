package units

import (
	"math"
	"strings"
	"testing"
)

// Reject non-finite: either caught at option boundary (WithMana) or by Validate() after options.
func TestValidateRejectsNonFinite(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		opt  Option
	}{
		// Offense path → caught by Validate()
		{"AD_NaN", WithAD(math.NaN())},
		{"AD_PosInf", WithAD(math.Inf(+1))},
		{"AD_NegInf", WithAD(math.Inf(-1))},

		// Defense path → caught by Validate()
		{"Armor_NaN", WithArmor(math.NaN())},
		{"Armor_PosInf", WithArmor(math.Inf(+1))},
		{"Armor_NegInf", WithArmor(math.Inf(-1))},

		// Resource path → caught at option boundary (WithMana)
		{"Mana_NaN", WithMana(math.NaN(), math.NaN(), math.NaN(), math.NaN(), math.NaN())},
		{"Mana_PosInf", WithMana(math.Inf(+1), math.Inf(+1), math.Inf(+1), math.Inf(+1), math.Inf(+1))},
		{"Mana_NegInf", WithMana(math.Inf(-1), math.Inf(-1), math.Inf(-1), math.Inf(-1), math.Inf(-1))},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			_, err := NewStats(tc.opt)
			if err == nil {
				t.Fatalf("expected error for non-finite input (%s), got nil", tc.name)
			}
			if !strings.Contains(err.Error(), "non-finite") {
				t.Fatalf("error should mention 'non-finite'; got: %v", err)
			}
		})
	}
}

func TestValidate_RangeConstraint(t *testing.T) {
	t.Parallel()

	// Invalid: must be able to hit adjacent hex (range < 1 should fail)
	_, err := NewStats(WithRange(0.5))
	if err == nil {
		t.Fatalf("expected error when range < 1")
	}
	if !strings.Contains(strings.ToLower(err.Error()), "range") {
		t.Fatalf("error should mention range, got: %v", err)
	}

	// Valid boundary
	_, err = NewStats(WithRange(1))
	if err != nil {
		t.Fatalf("unexpected error at range == 1: %v", err)
	}

	// Valid > 1
	_, err = NewStats(WithRange(3))
	if err != nil {
		t.Fatalf("unexpected error at range == 3: %v", err)
	}
}

func TestWithResource_RejectsNonFinite(t *testing.T) {
	t.Parallel()

	_, err := NewStats(
		WithRange(1),
		WithResource(Resource{
			ManaMin:   math.NaN(),
			ManaMax:   100,
			ManaStart: 0,
			ManaRegen: 5,
		}),
	)
	if err == nil {
		t.Fatalf("expected error for non-finite resource")
	}
	if !strings.Contains(err.Error(), "non-finite") {
		t.Fatalf("error should mention 'non-finite'; got: %v", err)
	}
}
