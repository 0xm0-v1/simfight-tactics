package units

import "math"

// sanitizeOffense enforces data invariants only.
// Mechanics (AS cap 5.0/6.0, crit overflow 2:1, anti-heal, etc.) live in the engine.
func sanitizeOffense(o OffenseStats) OffenseStats {
	return OffenseStats{
		Range:  nonNeg(o.Range),
		BaseAD: nonNeg(o.BaseAD),
		AD:     nonNeg(o.AD),
		AP:     nonNeg(o.AP),
		AS:     nonNeg(o.AS),
		// Do NOT clamp to [0,1] here: engine converts overflow >100% to crit damage (2:1).
		CritChance: nonNeg(o.CritChance),
		CritDamage: maxf(1.0, o.CritDamage),
		Omnivamp:   nonNeg(o.Omnivamp),
		DamageAmp:  o.DamageAmp, // may be negative or positive
	}
}

func sanitizeDefense(d DefenseStats) DefenseStats {
	return DefenseStats{
		HP:         nonNeg(d.HP),
		Armor:      nonNeg(d.Armor),
		MR:         nonNeg(d.MR),
		Durability: nonNeg(d.Durability),
	}
}

// sanitizeResource enforces cross-field invariants for Resource.
func sanitizeResource(r Resource) Resource {
	min := maxf(0, r.ManaMin)
	max := maxf(min, r.ManaMax)           // ensure max >= min
	start := clamp(r.ManaStart, min, max) // ensure start in [min, max]
	return Resource{
		ManaMin:   min,
		ManaMax:   max,
		ManaStart: start,
		ManaRegen: maxf(0, r.ManaRegen),
	}
}

// -----------------------------
// Math / guard utilities
// -----------------------------

// nonNeg clamps negatives to 0, but **does not** coerce NaN/±Inf;
// non-finites are left as-is so Validate() can fail fast.
func nonNeg(v float64) float64 {
	if math.IsNaN(v) || math.IsInf(v, 0) {
		return v // propagate non-finite to validation
	}
	if v < 0 {
		return 0
	}
	return v
}

func clamp(v, lo, hi float64) float64 {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

func maxf(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}
