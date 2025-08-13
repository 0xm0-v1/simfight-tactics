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
	Range      float64 `json:"range"`
	BaseAD     float64 `json:"base_attack_damage"`
	AD         float64 `json:"attack_damage"`
	AP         float64 `json:"ability_power"`
	AS         float64 `json:"attack_speed"`
	CritChance float64 `json:"critical_strike_chance"`
	CritDamage float64 `json:"critical_strike_damage"`
	Omnivamp   float64 `json:"omnivamp"`
	DamageAmp  float64 `json:"damage_amp"`
}

// DefenseStats represents defensive stats.
type DefenseStats struct {
	HP         float64 `json:"hp"`
	Armor      float64 `json:"armor"`
	MR         float64 `json:"magic_resist"`
	Durability float64 `json:"durability"`
}

// Resource represents resource-related values (e.g., mana).
type Resource struct {
	ManaMin   float64 `json:"mana_min"`
	ManaMax   float64 `json:"mana_max"`
	ManaStart float64 `json:"mana_start"`
	ManaRegen float64 `json:"mana_regen"`
}

func (s Stats) String() string {
	// Pretty-print via JSON to leverage json tags
	b, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Sprintf("Stats<marshal error: %v>", err)
	}
	return string(b)
}
