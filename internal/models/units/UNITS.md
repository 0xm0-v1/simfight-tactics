# Units Package - Implementation Guide

## Overview

This guide describes the standardized processes for implementing
statistics, role overrides, and unit creation in the Units package
within the simfight-tactics domain layer.

## ⚡ Quick Reference

| Task                   | How to do it | Key Files |
|------------------------|--------------|-----------|
| ➕ Add new stat field  | Edit model → add default → sanitize → option → validate | `stats.go`, `stats_default.go`, `stats_sanitize.go`, `stats_options.go`, `stats_validate.go` |
| 🎭 Configure role defaults | Define in `roles.json` | `roles.json` |
| 🏗️ Build a champion      | Use `BuildUnit(...)` with role + options | `unit_factory.go` |
| ✅ Validate stats        | Fail-fast checks for NaN / ±Inf | `stats_validate.go` |
| 🧪 Test stat behavior    | Unit + integration tests | `*_test.go` |


## Part 1: Stats Implementation

1.  **Model Definition**

``` go
// File: internal/domain/units/stats.go
type OffenseStats struct {
    AD float64 `json:"ad"`
    NewStat float64 `json:"new_stat"`
}
```

2.  **Default Values**

``` go
// File: internal/domain/units/stats_default.go
var defaultStats = Stats{
    Offense: OffenseStats{
        AD: 0,
        NewStat: 0, // baseline value
    },
}
```

3.  **Sanitization**

``` go
// File: internal/domain/units/stats_sanitize.go
func sanitizeOffense(o OffenseStats) OffenseStats {
    return OffenseStats{
        AD: nonNeg(o.AD),
        NewStat: nonNeg(o.NewStat),
    }
}
```

4.  **Functional Options**

``` go
// File: internal/domain/units/stats_options.go
func WithNewStat(v float64) Option {
    return setOffense(func(o *OffenseStats) { o.NewStat = v })
}
```

5.  **Validation**

``` go
// File: internal/domain/units/stats_validate.go
func (s Stats) Validate() error {
    if anyNonFinite(s.Offense.AD, s.Offense.NewStat) {
        return fmt.Errorf("stats contain non-finite values")
    }
    return nil
}
```

### Sanitization Helpers

-   **nonNeg(v)** → clamp to ≥ 0\
-   **clamp(v, lo, hi)** → bound to range\
-   **anyNaN(vals...)** → guard against NaN\
-   **maxf(a, b)** → primitive max

------------------------------------------------------------------------

## Part 2: Role Overrides

### Steps

1.  **Configure Role-Specific Overrides**

``` json
// File: internal/config/set15/roles.json
{
  "stats_per_roles": {
    "tank": {
      "defense": {
        "NewOverrideStats": 0.15
      }
    }
  }
}
```
------------------------------------------------------------------------

## Part 3: Initialization / Wiring

1.  **Validate Configuration**

``` go
// File: cmd/stfd/main.go
cfg, err := units.LoadRoles("roles.json")
if err != nil { return err }
if err := units.ValidateRolesConfig(cfg); err != nil {
    return err
}
```

2.  **Apply Overrides at Runtime**

``` go
stats, err := units.StatsForRole("Attack Tank", cfg)
```

------------------------------------------------------------------------

## Part 4: Unit Factory

### Example: Simple Champion

``` go
// File: cmd/stfd/main.go
garen, err := units.BuildUnit(
    "Garen",
    1,
    []string{"Warlord", "Vanguard"},
    []string{"Attack Tank"},
    cfg,
    units.WithHP(650),
    units.WithArmor(35),
    units.WithMR(35),
)
```

### Key Behaviors

-   First role in the slice is **primary role**
-   Validation ensures invariants
-   Invalid configs fail early

------------------------------------------------------------------------

## Complete Example: Adding Shield Mechanic

### Step 1: Add Shield to Stats Model

``` go
// File: internal/models/units/stats.go
type DefenseStats struct {
    HP         float64 `json:"hp"`
    Armor      float64 `json:"armor"`
    MR         float64 `json:"magic_resist"`
    Durability float64 `json:"durability"`
    Shield     float64 `json:"shield"` // NEW
}
```

### Step 2: Default Value

``` go
// File: internal/models/units/stats_default.go
Defense: DefenseStats{
    HP: 0, Armor: 0, MR: 0, Durability: 0,
    Shield: 0, // NEW
}
```

### Step 3: Sanitization

``` go
// File: internal/models/units/stats_sanitize.go
func sanitizeDefense(d DefenseStats) DefenseStats {
    return DefenseStats{
        HP: nonNeg(d.HP),
        Armor: nonNeg(d.Armor),
        MR: nonNeg(d.MR),
        Durability: nonNeg(d.Durability),
        Shield: nonNeg(d.Shield), // NEW
    }
}
```

### Step 4: Option

``` go
// File: internal/models/units/stats_options.go
func WithShield(v float64) Option {
    return setDefense(func(d *DefenseStats) { d.Shield = v })
}
```

### Step 5: Validation

``` go
// File: internal/models/units/stats_validate.go
if anyNonFinite(s.Defense.HP, s.Defense.Armor, s.Defense.MR, s.Defense.Shield) {
    return fmt.Errorf("defense contains non-finite values")
}
```

### Step 6: Role Override

``` json
// File: internal/config/set15/roles.json
{
  "stats_per_roles": {
    "tank": {
      "defense": {
        "shield": 150.0
      }
    }
  }
}
```

### Step 7: Unit Creation

``` go
// File: cmd/stfd/main.go
garen, err := units.BuildUnit(
    "Garen",
    1,
    []string{"Warlord", "Vanguard"},
    []string{"Attack Tank"},
    cfg,
    units.WithShield(200),
)
```

### Step 8: Tests

``` go
// File: internal/models/units/stats_shield_test.go
func TestShield_Sanitization(t *testing.T) {
    s := build(t, WithShield(-100))
    if s.Defense.Shield != 0 {
        t.Fatalf("negative shield should clamp to 0, got %v", s.Defense.Shield)
    }
}
```
