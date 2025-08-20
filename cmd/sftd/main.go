package main

import (
	"fmt"

	"github.com/0xm0-v1/simfight-tactics/internal/models/units"
)

const rolesPath = "internal/config/set15/roles.json"

func main() {
	// 1) Load config Role (Source of truth)
	cfg, err := units.LoadRoles(rolesPath)
	if err != nil {
		panic(fmt.Errorf("failed to load roles from %s: %w", rolesPath, err))
	}
	// 1.1) Activate Strict Mode
	cfg.Strict = true

	// 2) Build an Unit
	u, err := units.BuildUnit(
		"Garen",
		1,                               // Cost
		[]string{"Warlord", "Vanguard"}, // Traits
		[]string{"Attack Tank"},         // Roles
		cfg,
		// Explicit Overrides
		units.WithHP(650),
		units.WithArmor(35),
		units.WithMR(35),
		units.WithAD(55),
		units.WithAS(0.55),
		units.WithRange(1),
		// mana(min, max, start, regen/s, perHit)
		units.WithMana(0, 70, 30, 0, 10),
		units.WithCritChance(0.25),
		units.WithCritDamage(1.4),
	)
	if err != nil {
		panic(fmt.Errorf("failed to build unit: %w", err))
	}

	// 3) Display
	printUnit(&u)
}

func printUnit(u *units.Unit) {
	fmt.Printf("\nUnit created successfully!\n")
	fmt.Printf("- Name: %s\n", u.Name)
	fmt.Printf("- Cost: %d\n", u.Cost)
	fmt.Printf("- Traits: %v\n", u.Traits)
	fmt.Printf("- Roles: %v\n", u.Roles)

	fmt.Printf("\nBase Stats:\n")
	fmt.Printf("- HP: %.0f\n", u.Stats.Defense.HP)
	fmt.Printf("- Armor: %.0f\n", u.Stats.Defense.Armor)
	fmt.Printf("- Magic Resist: %.0f\n", u.Stats.Defense.MR)

	fmt.Printf("- Attack Damage: %.0f\n", u.Stats.Offense.AD)
	fmt.Printf("- Ability Power: %.0f\n", u.Stats.Offense.AP)
	fmt.Printf("- Attack Speed: %.2f\n", u.Stats.Offense.AS)
	fmt.Printf("- Range: %.0f\n", u.Stats.Offense.Range)

	fmt.Printf("- Crit Chance: %.0f%%\n", u.Stats.Offense.CritChance*100)
	fmt.Printf("- Crit Damage: %.0fx\n", u.Stats.Offense.CritDamage)

	fmt.Printf("- Mana: %.0f/%.0f (start: %.0f, regen: %.1f/s, per-hit: %.0f)\n",
		u.Stats.Resource.ManaMin,
		u.Stats.Resource.ManaMax,
		u.Stats.Resource.ManaStart,
		u.Stats.Resource.ManaRegen,
		u.Stats.Resource.ManaPerHit,
	)
}
