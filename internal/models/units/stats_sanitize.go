package units

import "math"

// sanitizeOffense enforces data invariants only.
func sanitizeOffense(o OffenseStats) OffenseStats {
	return OffenseStats{
		Range:      nonNeg(o.Range),
		BaseAD:     nonNeg(o.BaseAD),
		AD:         nonNeg(o.AD),
		AP:         nonNeg(o.AP),
		AS:         nonNeg(o.AS),
		CritChance: nonNeg(o.CritChance),    // engine handles >1 overflow
		CritDamage: maxf(1.0, o.CritDamage), // never below 1.0
		Omnivamp:   sanitizeOmnivamp(o.Omnivamp),
		DamageAmp:  o.DamageAmp, // may be negative or positive
	}
}

func sanitizeOmnivamp(v Omnivamp) Omnivamp {
	min := nonNeg(v.OmnivampMin)
	max := nonNeg(v.OmnivampMax)
	if max < min {
		max = min // ensure max >= min
	}
	// keep NaN/Inf as-is so Validate() can fail fast; else clamp to [min,max]
	current := v.CurrentOmnivamp
	if !(math.IsNaN(current) || math.IsInf(current, 0)) {
		current = clamp(current, min, max)
	}
	return Omnivamp{
		OmnivampMin:     min,
		OmnivampMax:     max,
		CurrentOmnivamp: current,
	}
}

func sanitizeDefense(d DefenseStats) DefenseStats {
	minTP := -1
	maxTP := 1
	return DefenseStats{
		HP:             nonNeg(d.HP),
		Armor:          nonNeg(d.Armor),
		MR:             nonNeg(d.MR),
		Durability:     nonNeg(d.Durability),
		TargetPriority: clamp(d.TargetPriority, float64(minTP), float64(maxTP)),
	}
}

func sanitizeManaFromDamage(m ManaFromDamage) ManaFromDamage {
	// Non-finites are propagated to validation by nonNeg
	pre := nonNeg(m.PreMitigationRatio)
	post := nonNeg(m.PostMitigationRatio)
	cap := nonNeg(m.PerInstanceCap)
	return ManaFromDamage{
		Enabled:             m.Enabled,
		PreMitigationRatio:  pre,
		PostMitigationRatio: post,
		PerInstanceCap:      cap,
	}
}

// sanitizeResource enforces cross-field invariants for Resource.
func sanitizeResource(r Resource) Resource {
	min := maxf(0, r.ManaMin)
	max := maxf(min, r.ManaMax)           // ensure max >= min
	start := clamp(r.ManaStart, min, max) // ensure start in [min, max]
	return Resource{
		ManaMin:        min,
		ManaMax:        max,
		ManaStart:      start,
		ManaRegen:      maxf(0, r.ManaRegen),
		ManaFromDamage: sanitizeManaFromDamage(r.ManaFromDamage),
		ManaPerHit:     nonNeg(r.ManaPerHit),
	}
}

// -----------------------------
// Math / guard utilities
// -----------------------------

// nonNeg clamps negatives to 0, but **does not** coerce NaN/±Inf;
// non-finites are left as-is so Validate() can fail fast.
func nonNeg(v float64) float64 {
	if math.IsNaN(v) || math.IsInf(v, 0) {
		return v
	}
	if v < 0 {
		return 0
	}
	return v
}

func clamp(v float64, lo, hi float64) float64 {
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
