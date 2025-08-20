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

func validateRange(valueRange float64) error {
	const minimumRange = 1.0
	if valueRange < minimumRange {
		return fmt.Errorf("range must be >= %v, (got %v)", minimumRange, valueRange)
	}
	return nil
}

// Validate performs fail-fast checks we do NOT want to auto-correct.
func (s Stats) Validate() error {
	// 1) Non-finite numbers are always invalid.
	if anyNonFinite(
		// Offense
		s.Offense.Range, s.Offense.BaseAD, s.Offense.AD, s.Offense.AP, s.Offense.AS,
		s.Offense.CritChance, s.Offense.CritDamage, s.Offense.Omnivamp.OmnivampMax, s.Offense.Omnivamp.OmnivampMin, s.Offense.Omnivamp.CurrentOmnivamp, s.Offense.DamageAmp,
		// Defense
		s.Defense.HP, s.Defense.Armor, s.Defense.MR, s.Defense.Durability, s.Defense.TargetPriority,
		// Resource (flat)
		s.Resource.ManaMin, s.Resource.ManaMax, s.Resource.ManaStart, s.Resource.ManaRegen, s.Resource.ManaPerHit,
		// Resource (mana from damage)
		s.Resource.ManaFromDamage.PreMitigationRatio,
		s.Resource.ManaFromDamage.PostMitigationRatio,
		s.Resource.ManaFromDamage.PerInstanceCap,
	) {
		return fmt.Errorf("stats contain non-finite values (NaN/±Inf)")
	}

	// 2) Range lower bound (we don't auto-correct this one).
	if err := validateRange(s.Offense.Range); err != nil {
		return err
	}

	return nil
}
