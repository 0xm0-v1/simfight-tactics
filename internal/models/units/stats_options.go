package units

import "fmt"

// Option is the functional option type for Stats.
type Option func(*Stats) error

// NewStats builds Stats from Default() + options, validates, then sanitizes.
func NewStats(opts ...Option) (Stats, error) {
	s := Default()
	for _, opt := range opts {
		if err := opt(&s); err != nil {
			return Stats{}, err
		}
	}
	// Fail-fast validation BEFORE final normalization
	if err := s.Validate(); err != nil {
		return Stats{}, err
	}
	return s.normalized(), nil
}

// With applies options on a copy, validates, then sanitizes.
func (s Stats) With(opts ...Option) (Stats, error) {
	cp := s
	for _, opt := range opts {
		if err := opt(&cp); err != nil {
			return Stats{}, err
		}
	}
	if err := cp.Validate(); err != nil {
		return Stats{}, err
	}
	return cp.normalized(), nil
}

// --- Offense (granular setters) ---
func WithRange(v float64) Option      { return setOffense(func(o *OffenseStats) { o.Range = v }) }
func WithBaseAD(v float64) Option     { return setOffense(func(o *OffenseStats) { o.BaseAD = v }) }
func WithAD(v float64) Option         { return setOffense(func(o *OffenseStats) { o.AD = v }) }
func WithAP(v float64) Option         { return setOffense(func(o *OffenseStats) { o.AP = v }) }
func WithAS(v float64) Option         { return setOffense(func(o *OffenseStats) { o.AS = v }) }
func WithCritChance(v float64) Option { return setOffense(func(o *OffenseStats) { o.CritChance = v }) }
func WithCritDamage(v float64) Option { return setOffense(func(o *OffenseStats) { o.CritDamage = v }) }
func WithOmnivamp(v Omnivamp) Option  { return setOffense(func(o *OffenseStats) { o.Omnivamp = v }) }
func WithDamageAmp(v float64) Option  { return setOffense(func(o *OffenseStats) { o.DamageAmp = v }) }

func WithOmnivampValues(min, max, current float64) Option {
	return WithOmnivamp(Omnivamp{
		OmnivampMin:     min,
		OmnivampMax:     max,
		CurrentOmnivamp: current,
	})
}

// Bulk setter with sanitization.
func WithOffense(off OffenseStats) Option {
	return func(s *Stats) error {
		s.Offense = sanitizeOffense(off)
		return nil
	}
}

// --- Defense ---
func WithHP(v float64) Option         { return setDefense(func(d *DefenseStats) { d.HP = v }) }
func WithArmor(v float64) Option      { return setDefense(func(d *DefenseStats) { d.Armor = v }) }
func WithMR(v float64) Option         { return setDefense(func(d *DefenseStats) { d.MR = v }) }
func WithDurability(v float64) Option { return setDefense(func(d *DefenseStats) { d.Durability = v }) }
func WithTargetPriority(v float64) Option {
	return setDefense(func(d *DefenseStats) { d.TargetPriority = v })
}

// Bulk setter with sanitization.
func WithDefense(def DefenseStats) Option {
	return func(s *Stats) error {
		s.Defense = sanitizeDefense(def)
		return nil
	}
}

// --- Resource ---
func WithMana(min, max, start, regen, perHit float64) Option {
	return func(s *Stats) error {
		if anyNonFinite(min, max, start, regen, perHit) {
			return fmt.Errorf("resource contains non-finite values")
		}
		s.Resource = sanitizeResource(Resource{
			ManaMin:        min,
			ManaMax:        max,
			ManaStart:      start,
			ManaRegen:      regen,
			ManaFromDamage: s.Resource.ManaFromDamage, // preserve current sub-struct
			ManaPerHit:     perHit,
		})
		return nil
	}
}

func WithResource(res Resource) Option {
	return func(s *Stats) error {
		if anyNonFinite(
			res.ManaMin, res.ManaMax, res.ManaStart, res.ManaRegen, res.ManaPerHit,
			res.ManaFromDamage.PreMitigationRatio,
			res.ManaFromDamage.PostMitigationRatio,
			res.ManaFromDamage.PerInstanceCap,
		) {
			return fmt.Errorf("resource contains non-finite values")
		}
		s.Resource = sanitizeResource(res)
		return nil
	}
}

// WithManaFromDamage sets the mana-from-hit rules.
func WithManaFromDamage(enabled bool, preRatio, postRatio, perHitCap float64) Option {
	return func(s *Stats) error {
		if anyNonFinite(preRatio, postRatio, perHitCap) {
			return fmt.Errorf("mana_from_damage contains non-finite values")
		}
		cur := s.Resource
		cur.ManaFromDamage = sanitizeManaFromDamage(ManaFromDamage{
			Enabled:             enabled,
			PreMitigationRatio:  preRatio,
			PostMitigationRatio: postRatio,
			PerInstanceCap:      perHitCap,
		})
		s.Resource = sanitizeResource(cur)
		return nil
	}
}

// -----------------------------
// Internals
// -----------------------------

func setOffense(f func(*OffenseStats)) Option {
	return func(s *Stats) error {
		if s == nil {
			return fmt.Errorf("nil Stats")
		}
		tmp := s.Offense
		f(&tmp)                          // set RAW input
		s.Offense = sanitizeOffense(tmp) // then sanitize deterministically
		return nil
	}
}

func setDefense(f func(*DefenseStats)) Option {
	return func(s *Stats) error {
		if s == nil {
			return fmt.Errorf("nil Stats")
		}
		tmp := s.Defense
		f(&tmp)
		s.Defense = sanitizeDefense(tmp)
		return nil
	}
}

// normalized applies final sanitation only (data invariants).
// All fail-fast validation happens BEFORE calling this.
func (s Stats) normalized() Stats {
	s.Offense = sanitizeOffense(s.Offense)
	s.Defense = sanitizeDefense(s.Defense)
	s.Resource = sanitizeResource(s.Resource)
	return s
}
