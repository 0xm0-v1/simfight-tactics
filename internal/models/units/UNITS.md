# Units Package

1) **Model** – add the field with a clear JSON tag.
```go
// File: simfight-tactics/internal/domain/units/stats.go
type OffenseStats struct {
    // ...
    NewStats float64 `json:"new"`
}
```
2) **Defaults** – set a safe baseline.
```go
// File: simfight-tactics/internal/domain/units/stats_defaults.go
var defaultStats = Stats{
    Offense: OffenseStats{
        // ...
        NewStats: 0, // baseline
    },
    // ...
}
```
3) **Sanitize** (data invariants only) – no game logic here.
```go
// File: simfight-tactics/internal/domain/units/stats_sanitize.go
func sanitizeOffense(o OffenseStats) OffenseStats {
    return OffenseStats{
        // ...
        NewStats: maxf(0, o.NewStats), // non-negative
    }
}
```
4) **Functional Option** – ergonomic setter (+ covered by bulk setters).
```go
// File: simfight-tactics/internal/domain/units/stats_options.go
// Granular setter
func WithNewStats(v float64) Option {
    return setOffense(func(o *OffenseStats) { o.NewStats = maxf(0, v) })
}
```

### Helper cheat sheet (`stats_sanitize.go`) — compact, human-friendly

**Rule:** enforce **data invariants only** here. **Gameplay caps/conversions live in the engine.**

---

**`anyNaN(vals...) bool` — reject non-finite input**
- **Use:** upfront guard in options / `normalized()` → `if anyNaN(...) { return err }`
- **Avoid:** coercion (never “fix” NaN here).
- **Example:** `if anyNaN(min, max, start, regen) return fmt.Errorf("mana contains NaN")`

**`nonNeg(v) float64` — clamp to ≥ 0 (propagates NaN)**
- **Use:** AD, AP, AS, HP, Armor, MR, Durability, Omnivamp, Mana*.
- **Avoid:** fields that may be negative (e.g., `DamageAmp`).
- **Example:** `o.AD = nonNeg(o.AD)`  // if NaN → stays NaN and will be caught by validation

**`clamp(v, lo, hi) float64` — bound to [lo, hi]**
- **Use:** structural ranges (e.g., `ManaStart ∈ [ManaMin, ManaMax]`).
- **Avoid:** gameplay limits (no crit 100% cap, no AS 5/6 cap here).
- **Example:** `start = clamp(r.ManaStart, min, max)`

**`maxf(a, b) float64` — primitive max (does not handle NaN)**
- **Use:** compose invariants (`maxf(0, v)`, `maxf(min, max)`).
- **Avoid:** direct use with potentially NaN inputs; guard with `anyNaN` first.
- **Example:** `max = maxf(min, r.ManaMax)`  // ensures `max >= min`

---

**One-liners you’ll actually write**
```go
o.CritChance = nonNeg(o.CritChance) // no upper clamp here; engine handles overflow
d.Armor      = nonNeg(d.Armor)
r.ManaStart  = clamp(r.ManaStart, maxf(0, r.ManaMin), maxf(maxf(0, r.ManaMin), r.ManaMax))
