# Units Package - Stats Implementation Guide

## Overview

This guide describes the standardized process for adding new statistics to the Units package in the simfight-tactics domain layer.

## Implementation Steps

### 1. Model Definition

Add the new field to the appropriate stats structure with a clear JSON tag.

```go
// File: simfight-tactics/internal/domain/units/stats.go

type OffenseStats struct {
    // ...existing fields
    NewStat float64 `json:"new_stat"`
}
```

### 2. Default Values

Set a safe baseline value in the defaults configuration.

```go
// File: simfight-tactics/internal/domain/units/stats_default.go

var defaultStats = Stats{
    Offense: OffenseStats{
        // ...existing defaults
        NewStat: 0, // baseline value
    },
    // ...
}
```

### 3. Sanitization

Apply data invariants only. Game logic belongs in the engine layer.

```go
// File: simfight-tactics/internal/domain/units/stats_sanitize.go

func sanitizeOffense(o OffenseStats) OffenseStats {
    return OffenseStats{
        // ...existing sanitization
        NewStat: nonNeg(o.NewStat), // clamp ≥ 0, propagate NaN/±Inf
    }
}
```

### 4. Functional Options

Create an ergonomic setter using the Option pattern.

```go
// File: simfight-tactics/internal/domain/units/stats_options.go

func WithNewStat(v float64) Option {
    return setOffense(func(o *OffenseStats) { 
        o.NewStat = v // raw value; sanitized later
    })
}
```

### 5. Validation

Add fail-fast checks for non-finite values.

```go
// File: simfight-tactics/internal/domain/units/stats_validate.go

func (s Stats) Validate() error {
    if anyNonFinite(
        // ...existing fields
        s.Offense.NewStat, // add if relevant
    ) {
        return fmt.Errorf("stats contain non-finite values (NaN/±Inf)")
    }
    // other fail-fast rules...
    return nil
}
```

## Special Case: Resource Fields

When adding fields to the Resource structure, follow these additional steps:

1. **Add to defaults**
   ```go
   defaultStats.Resource.NewResourceStat = 0
   ```

2. **Sanitize in dedicated function**
   ```go
   func sanitizeResource() {
       // Apply appropriate sanitization
   }
   ```

3. **Include in validation guard**
   ```go
   if anyNonFinite(
       res.ManaMin, res.ManaMax, res.ManaStart, res.ManaRegen, res.ManaPerHit,
       res.ManaFromDamage.PreMitigationRatio,
       res.ManaFromDamage.PostMitigationRatio,
       res.ManaFromDamage.PerInstanceCap,
       res.NewResourceStat, // new field
   ) {
       return fmt.Errorf("resource contains non-finite values")
   }
   ```

## Sanitization Helper Reference

### Core Principle
Enforce **data invariants only**. Gameplay caps and conversions belong in the engine layer.

### Helper Functions

#### `anyNaN(vals...) bool`
**Purpose:** Reject non-finite input  
**Usage:** Upfront guard in options/normalized  
**Avoid:** Coercion (never "fix" NaN)  
**Example:**
```go
if anyNaN(min, max, start, regen) {
    return fmt.Errorf("mana contains NaN")
}
```

#### `nonNeg(v) float64`
**Purpose:** Clamp to ≥ 0 (propagates NaN)  
**Usage:** AD, AP, AS, HP, Armor, MR, Durability, Omnivamp, Mana*  
**Avoid:** Fields that may be negative (e.g., DamageAmp)  
**Example:**
```go
o.AD = nonNeg(o.AD)
```

#### `clamp(v, lo, hi) float64`
**Purpose:** Bound value to [lo, hi]  
**Usage:** Structural ranges (e.g., ManaStart ∈ [ManaMin, ManaMax])  
**Avoid:** Gameplay limits (no crit 100% cap, no AS 5/6 cap)  
**Example:**
```go
start = clamp(r.ManaStart, min, max)
```

#### `maxf(a, b) float64`
**Purpose:** Primitive max (no NaN handling)  
**Usage:** Compose invariants  
**Avoid:** Direct use with potentially NaN inputs  
**Example:**
```go
max = maxf(min, r.ManaMax)
```

### Common Patterns

```go
o.CritChance = nonNeg(o.CritChance)
d.Armor = nonNeg(d.Armor)
r.ManaStart = clamp(
    r.ManaStart, 
    maxf(0, r.ManaMin), 
    maxf(maxf(0, r.ManaMin), r.ManaMax)
)
```

---

# Units Package — Stats & Role Overrides Guide

This document explains how **Role-based overrides** work in the Units domain: how to define them in JSON, load/validate them, and apply them deterministically to `Stats`.

> TL;DR  
> - **Single source of truth**: `roles.json`  
> - **Deterministic**: fixed parsing & validation  
> - **Transparent**: unknown keys & type mismatches are reported

---

## 1) Glossary

- **Damage Type tokens**: `attack`, `magic`, `hybrid` (configurable).
- **Role keys**: `tank`, `fighter`, `assassin`, `marksman`, `caster`, `specialist` (configurable).
- **Role label**: free-form label that mixes damage type + role key, e.g. `"Attack Tank"`, `"magic-caster"`.  
  It is normalized to `(roleKey, damageType)`; the **roleKey** drives overrides today.

---

## 2) roles.json — Structure & Example

The JSON lives next to your binaries/configs (path is app-specific). Example:

```json
{
  "sft": {
    "version": "set15-15.2",
    "updated_at": "2025-08-12",
    "sources": ["...patch links..."]
  },

  "damage_type": ["Attack", "Magic", "Hybrid"],
  "role_type":   ["Tank", "Fighter", "Assassin", "Marksman", "Caster", "Specialist"],

  "roles": [
    "Attack Tank", "Attack Fighter", "...", "Hybrid Specialist"
  ],

  "stats_per_roles": {
    "tank": {
      "defense": { "target_priority": 1.0 },
      "resource": {
        "mana_from_damage": {
          "enabled": true,
          "pre_mitigation_ratio": 0.01,
          "post_mitigation_ratio": 0.03,
          "per_instance_cap": 42.5
        },
        "mana_per_hit": 5.0
      }
    },
    "fighter": {
      "offense": {
        "omnivamp": {
          "omnivamp_min": 0.08,
          "omnivamp_max": 0.2,
          "current_omnivamp": 0.08
        }
      },
      "resource": { "mana_per_hit": 10.0 }
    },
    "assassin": {
      "defense": { "target_priority": -1.0 },
      "resource": { "mana_per_hit": 10.0 }
    },
    "marksman": { "resource": { "mana_per_hit": 10.0 } },
    "caster":   { "resource": { "mana_per_hit": 7.0, "mana_regen": 2.0 } }
  }
}
```

### Mapping rule
`stats_per_roles` mirrors the Go `Stats` structure **by JSON tags**. Nested objects must match the nested structs:
- `defense.target_priority` → `Stats.Defense.TargetPriority` (`float64`)
- `resource.mana_per_hit` → `Stats.Resource.ManaPerHit` (`float64`)
- `resource.mana_from_damage.*` → nested `ManaFromDamage` struct
- `offense.omnivamp.*` → nested `Omnivamp` struct

Unknown keys or type mismatches are reported (see §5).

---

## 3) Loading & Normalization

### Load the configuration
```go
cfg, err := units.LoadRoles(pathToRolesJSON)
if err != nil { /* handle */ }
```

- `cfg.RoleTypes` and `cfg.DamageTypes` provide **allowed tokens**.  
  If missing in JSON, the loader falls back to built-ins:
  - Roles: `tank|fighter|assassin|marksman|caster|specialist`
  - Damage: `attack|magic|hybrid`

### Role label → (roleKey, damageType)
```go
roleKey, dmgType, ok := detectRoleKey("Attack Tank", cfg.ValidRoleKeys(), cfg.DamageTypeTokens())
```

Rules:
- Case-insensitive; splits on space, tab, `-`, `_`, `/`.
- **Exactly one** role token must be present.
- At most one valid damage token can be present (optional).
- Any unknown token **before** the role token ⇒ invalid label.

---

## 4) Applying Role overrides to Stats

### Getting stats for a role
```go
stats, err := units.StatsForRole("Attack Tank", cfg)
```

Flow:
1. Build base stats via `NewStats(...)` (defaults + options).
2. Normalize the label → `roleKey`.
3. Look up `cfg.StatsPerRoles[roleKey]`.
4. Apply the JSON object onto `Stats` using JSON tags.
5. Validate/sanitize the result.

### Creating a Unit with role-driven stats
```go
u, err := units.BuildUnit("Garen", 1, []string{"Warlord"}, []string{"Attack Tank"}, cfg)
```

- The **first** role in `roles[]` is the **primary** and drives overrides today.

---

## 5) Validation & Strictness

### Config-time validation
```go
if err := units.ValidateRolesConfig(cfg); err != nil {
    // Unknown roles, non-object entries, unknown keys, or type errors
}
```

### Apply-time reporting
- In `StatsForRole(...)`, the apply step collects:
  - `UnknownKeys`: keys that don’t match any JSON tag path
  - `TypeErrors`: wrong primitive types (expects `number` or `boolean` as appropriate)
- Behavior depends on `cfg.Strict`:
  - `false` (default): **log warnings**, continue.
  - `true`: **fail-fast** with error.

### Stats invariants
After overrides, `Stats.Validate()` + normalization enforce data invariants only (non-negativity, finite numbers, structural ranges). Gameplay caps are **not** enforced here.

---

## 6) Supported field types in overrides

- `float`-like numbers (JSON numbers; ints are accepted and coerced to `float64`)
- `boolean` for flags (e.g., `resource.mana_from_damage.enabled`)
- Nested objects must be JSON objects
- Slices/maps/strings are **not** supported in role overrides as of now (reported as type errors).

---

## 7) Examples

### A. Make Fighters lifesteal by default
```json
"fighter": {
  "offense": {
    "omnivamp": {
      "omnivamp_min": 0.08,
      "omnivamp_max": 0.20,
      "current_omnivamp": 0.08
    }
  },
  "resource": { "mana_per_hit": 10.0 }
}
```

### B. Tanks gain mana when hit
```json
"tank": {
  "defense": { "target_priority": 1.0 },
  "resource": {
    "mana_from_damage": {
      "enabled": true,
      "pre_mitigation_ratio": 0.01,
      "post_mitigation_ratio": 0.03,
      "per_instance_cap": 42.5
    },
    "mana_per_hit": 5.0
  }
}
```

### C. Strict mode for CI
```go
cfg, _ := units.LoadRoles("roles.json")
cfg.Strict = true
if err := units.ValidateRolesConfig(cfg); err != nil { return err } // fail the build
```

---

## 8) Testing Checklist

- **Determinism**: same seed + same `roles.json` ⇒ identical `Stats`.
- **Validation**: inject unknown keys / wrong types and assert:
  - `ValidateRolesConfig` catches them.
  - With `Strict=true`, `StatsForRole` returns error.
- **Normalization**: try labels like `" MAGIC /   CASTER "`, `"attack-marksman"`, `"hybrid_assassin"`; ensure they normalize.
- **Sanitization**: negative or NaN in JSON should be rejected by `Stats.Validate()` (fail-fast).

---

## 9) Troubleshooting

- `invalid role label`: the string didn’t contain a valid role token or had conflicting tokens.
- `override issues: unknown_keys=...`: key path doesn’t match any `json` tag in `Stats`; fix the JSON path or the struct tags.
- `type_errors=...`: supply a JSON number (or boolean) as expected by the target field.
- No overrides found for role: if `Strict=false`, it logs a notice and uses base stats.

---

## 10) Design Notes

- **Data-driven**: role/damage tokens come from JSON (with safe fallbacks).
- **JSON-tag reflection**: keeps the mapping transparent; changing struct tags updates the contract.
- **Domain purity**: engine-level gameplay rules (caps, conversions, RNG) remain out of `Stats`.