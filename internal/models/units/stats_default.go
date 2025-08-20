package units

// defaultStats is unexported to prevent accidental global mutation.
var defaultStats = Stats{
	Offense: OffenseStats{
		Range:      1,
		BaseAD:     0,
		AD:         0,
		AP:         0,
		AS:         0,
		CritChance: 0.25,
		CritDamage: 1.4, // baseline crit damage = 140%
		DamageAmp:  0,
		Omnivamp: Omnivamp{
			OmnivampMin:     0,
			OmnivampMax:     0,
			CurrentOmnivamp: 0,
		}, // Default for fighter is +8% - +20% based on Stage.
	},
	Defense: DefenseStats{
		HP:             0,
		Armor:          0,
		MR:             0,
		Durability:     0,
		TargetPriority: 0,
	},
	Resource: Resource{
		ManaMin:   0,
		ManaMax:   0,
		ManaStart: 0,
		ManaRegen: 0,
		ManaFromDamage: ManaFromDamage{
			Enabled:             false,
			PreMitigationRatio:  0,
			PostMitigationRatio: 0,
			PerInstanceCap:      0,
		},
		ManaPerHit: 0,
	},
}

// Default returns a copy of default values.
func Default() Stats { return defaultStats }
