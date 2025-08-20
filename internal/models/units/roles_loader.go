package units

import (
	"fmt"
	"os"
	"strings"

	json "encoding/json/v2"
)

// --- built-in fallbacks to keep robustness if JSON omits lists ---
var builtinValidRoleKeys = []string{"tank", "fighter", "assassin", "marksman", "caster", "specialist"}
var builtinDamageTypeTokens = []string{"attack", "magic", "hybrid"}

type RolesLoader struct {
	StatsPerRoles map[string]any `json:"stats_per_roles"`
	RoleTypes     []string       `json:"role_type"`   // NEW
	DamageTypes   []string       `json:"damage_type"` // NEW
	Strict        bool           `json:"-"`           // runtime-only
}

func LoadRoles(path string) (RolesLoader, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return RolesLoader{}, fmt.Errorf("read roles config: %w", err)
	}
	var cfg RolesLoader
	if err := json.Unmarshal(b, &cfg); err != nil {
		return RolesLoader{}, fmt.Errorf("parse roles config: %w", err)
	}
	cfg.Strict = false
	return cfg, nil
}

// Lowercases and returns a set for fast membership checks.
// Falls back to built-in lists if the JSON list is empty.
func (c RolesLoader) ValidRoleKeys() map[string]struct{} {
	src := c.RoleTypes
	if len(src) == 0 {
		src = builtinValidRoleKeys
	}
	out := make(map[string]struct{}, len(src))
	for _, s := range src {
		out[strings.ToLower(strings.TrimSpace(s))] = struct{}{}
	}
	return out
}

func (c RolesLoader) DamageTypeTokens() map[string]struct{} {
	src := c.DamageTypes
	if len(src) == 0 {
		src = builtinDamageTypeTokens
	}
	out := make(map[string]struct{}, len(src))
	for _, s := range src {
		out[strings.ToLower(strings.TrimSpace(s))] = struct{}{}
	}
	return out
}
