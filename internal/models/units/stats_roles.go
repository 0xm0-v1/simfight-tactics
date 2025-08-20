package units

import (
	"fmt"
	"log"
	"strings"
)

func findRoleOverrideMap(sp map[string]any, roleKey string) map[string]any {
	if sp == nil {
		return nil
	}
	// correspondance exacte
	if v, ok := sp[roleKey]; ok {
		if m, ok := v.(map[string]any); ok {
			return m
		}
	}
	// correspondance lowercase (sécurité si la config a des capitales)
	lc := strings.ToLower(roleKey)
	if v, ok := sp[lc]; ok {
		if m, ok := v.(map[string]any); ok {
			return m
		}
	}
	return nil
}

func StatsForRole(role string, cfg RolesLoader, opts ...Option) (Stats, error) {
	base, err := NewStats(opts...)
	if err != nil {
		return Stats{}, err
	}

	// NEW: derive allowed tokens from cfg (with fallbacks)
	validRoles := cfg.ValidRoleKeys()
	damageTokens := cfg.DamageTypeTokens()

	roleKey, dmgType, ok := detectRoleKey(role, validRoles, damageTokens)
	if !ok {
		return Stats{}, fmt.Errorf("invalid role label %q: must contain a valid role token and optional valid damage type", role)
	}
	_ = dmgType // still validated but not used for overrides yet

	roleMap := findRoleOverrideMap(cfg.StatsPerRoles, roleKey)
	if len(roleMap) == 0 {
		if !cfg.Strict {
			log.Printf("[roles] no overrides for role=%q (normalized=%q)", role, roleKey)
		}
		return base, nil
	}

	applied := base
	report := applyRoleMapToStats(&applied, roleMap)

	if cfg.Strict && !report.empty() {
		return Stats{}, fmt.Errorf("invalid role stats (%s): override issues: unknown_keys=%v, type_errors=%v",
			role, report.UnknownKeys, report.TypeErrors,
		)
	}
	if !cfg.Strict && (len(report.UnknownKeys) > 0 || len(report.TypeErrors) > 0) {
		log.Printf("[roles] role=%q override warnings: unknown_keys=%v, type_errors=%v",
			role, report.UnknownKeys, report.TypeErrors)
	}

	if err := applied.Validate(); err != nil {
		return Stats{}, fmt.Errorf("invalid role stats (%s): %w", role, err)
	}
	return applied.normalized(), nil
}
