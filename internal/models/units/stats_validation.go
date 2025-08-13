package units

import (
	"fmt"
	"math"
)

// anyNonFinite returns true if any value is NaN or ±Inf (fail-fast guard).
func anyNonFinite(vals ...float64) bool {
	for _, v := range vals {
		if math.IsNaN(v) || math.IsInf(v, 0) {
			return true
		}
	}
	return false
}

// Validate performs fail-fast checks we do NOT want to auto-correct.
// It does NOT mutate s. Sanitation (non-negativity, clamping) happens elsewhere.
func (s Stats) Validate() error {
	// 1) Non-finite numbers are always invalid.
	if anyNonFinite(
		s.Offense.Range, s.Offense.BaseAD, s.Offense.AD, s.Offense.AP, s.Offense.AS,
		s.Offense.CritChance, s.Offense.CritDamage, s.Offense.Omnivamp, s.Offense.DamageAmp,
		s.Defense.HP, s.Defense.Armor, s.Defense.MR, s.Defense.Durability,
		s.Resource.ManaMin, s.Resource.ManaMax, s.Resource.ManaStart, s.Resource.ManaRegen,
	) {
		return fmt.Errorf("stats contain non-finite values (NaN/±Inf)")
	}

	// 2) Champion constraint: must be able to hit adjacent hexes.
	if s.Offense.Range < 1 {
		return fmt.Errorf("range must be >= 1")
	}
	return nil
}
