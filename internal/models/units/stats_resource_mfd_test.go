package units

import (
	"math"
	"strings"
	"testing"
)

func TestWithManaFromDamage_NonFiniteRejected(t *testing.T) {
	t.Parallel()
	_, err := NewStats(WithRange(1), WithManaFromDamage(true, math.NaN(), 0.03, 40))
	if err == nil || !strings.Contains(err.Error(), "non-finite") {
		t.Fatalf("expected non-finite error, got %v", err)
	}
}

func TestWithManaFromDamage_SanitizeNonNegatives(t *testing.T) {
	t.Parallel()
	s := build(t, WithManaFromDamage(true, -0.01, 0.02, -5))
	mfd := s.Resource.ManaFromDamage
	if !mfd.Enabled || mfd.PreMitigationRatio != 0 || mfd.PostMitigationRatio != 0.02 || mfd.PerInstanceCap != 0 {
		t.Fatalf("unexpected MFD after sanitize: %+v", mfd)
	}
}

func TestWithResource_NonFiniteGuard_IncludingSubstruct(t *testing.T) {
	t.Parallel()
	_, err := NewStats(
		WithRange(1),
		WithResource(Resource{
			ManaMin: 0, ManaMax: 100, ManaStart: 0, ManaRegen: 5, ManaPerHit: 1,
			ManaFromDamage: ManaFromDamage{Enabled: true, PreMitigationRatio: math.NaN()},
		}),
	)
	if err == nil || !strings.Contains(err.Error(), "resource contains non-finite") {
		t.Fatalf("expected resource non-finite error, got %v", err)
	}
}
