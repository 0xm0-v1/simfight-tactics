package units

import (
	"encoding/json"
	"fmt"
)

// Stats groups offensive, defensive, and resource stats.
type Stats struct {
	Offense  OffenseStats `json:"offense"`
	Defense  DefenseStats `json:"defense"`
	Resource Resource     `json:"resource"`
}

// OffenseStats represents offensive stats.
type OffenseStats struct {
	Range      float64  `json:"range"`
	BaseAD     float64  `json:"base_attack_damage"`
	AD         float64  `json:"attack_damage"`
	AP         float64  `json:"ability_power"`
	AS         float64  `json:"attack_speed"`
	CritChance float64  `json:"critical_strike_chance"`
	CritDamage float64  `json:"critical_strike_damage"`
	Omnivamp   Omnivamp `json:"omnivamp"`
	DamageAmp  float64  `json:"damage_amp"`
}

type Omnivamp struct {
	OmnivampMin     float64 `json:"omnivamp_min"`
	OmnivampMax     float64 `json:"omnivamp_max"`
	CurrentOmnivamp float64 `json:"current_omnivamp"`
}

// DefenseStats represents defensive stats.
type DefenseStats struct {
	HP         float64 `json:"hp"`
	Armor      float64 `json:"armor"`
	MR         float64 `json:"magic_resist"`
	Durability float64 `json:"durability"`
	// TargetPriority is a tie-breaker bias in [-1,+1]:
	// -1 = less likely to be targeted, +1 = more likely, 0 = neutral.
	// Engine should only apply it when distance (and other primary criteria) are tied.
	TargetPriority float64 `json:"target_priority"`
}

// Resource represents resource-related values (e.g., mana).
type Resource struct {
	ManaMin        float64        `json:"mana_min"`
	ManaMax        float64        `json:"mana_max"`
	ManaStart      float64        `json:"mana_start"`
	ManaRegen      float64        `json:"mana_regen"`
	ManaFromDamage ManaFromDamage `json:"mana_from_damage"`
	ManaPerHit     float64        `json:"mana_per_hit"`
}

// ManaFromDamage defines how mana is generated when taking damage.
// Mana gain per instance=(0.01×damage_taken_pre_mitigation)+(0.03×damage_taken_post_mitigation)
type ManaFromDamage struct {
	Enabled             bool    `json:"enabled"`
	PreMitigationRatio  float64 `json:"pre_mitigation_ratio"`  // mana per damage point BEFORE reductions
	PostMitigationRatio float64 `json:"post_mitigation_ratio"` // mana per damage point AFTER reductions
	PerInstanceCap      float64 `json:"per_instance_cap"`      // max mana gained per hit instance (0 = no cap)
}

func (s Stats) String() string {
	// Pretty-print via JSON to leverage json tags
	b, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Sprintf("Stats<marshal error: %v>", err)
	}
	return string(b)
}
