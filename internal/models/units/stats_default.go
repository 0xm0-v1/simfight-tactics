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
		// Baseline crit damage = 140% (1.4)
		CritDamage: 1.4,
		DamageAmp:  0,
		Omnivamp:   0,
	},
	Defense: DefenseStats{
		HP:         0,
		Armor:      0,
		MR:         0,
		Durability: 0,
	},
	Resource: Resource{
		ManaMin:   0,
		ManaMax:   0,
		ManaStart: 0,
		ManaRegen: 0,
	},
}

// Default returns a copy of default values.
func Default() Stats { return defaultStats }
