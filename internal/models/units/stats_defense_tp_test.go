package units

import "testing"

func TestSanitize_TargetPriority_Clamp(t *testing.T) {
	t.Parallel()

	s := build(t, WithTargetPriority(2.5))
	if s.Defense.TargetPriority != 1 {
		t.Fatalf("TP should clamp to +1, got %v", s.Defense.TargetPriority)
	}
	s2 := build(t, WithTargetPriority(-3.0))
	if s2.Defense.TargetPriority != -1 {
		t.Fatalf("TP should clamp to -1, got %v", s2.Defense.TargetPriority)
	}
}
