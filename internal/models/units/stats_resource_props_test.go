package units

import "testing"

func TestSanitizeResource_Invariants(t *testing.T) {
	t.Parallel()

	// max < min, start > max, regen < 0, perHit < 0
	s := build(t, WithMana(80, 30, 150, -3, -10))
	r := s.Resource
	if !(r.ManaMin == 80 && r.ManaMax == 80 && r.ManaStart == 80) {
		t.Fatalf("expected min==max==start==80, got %+v", r)
	}
	if r.ManaRegen != 0 {
		t.Fatalf("ManaRegen clamps to 0, got %v", r.ManaRegen)
	}
	if r.ManaPerHit != 0 {
		t.Fatalf("ManaPerHit clamps to 0, got %v", r.ManaPerHit)
	}
}
