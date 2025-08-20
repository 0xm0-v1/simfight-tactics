package units

import (
	"fmt"
	"sort"
	"strings"
)

func ValidateRolesConfig(cfg RolesLoader) error {
	sp := cfg.StatsPerRoles
	if sp == nil {
		return nil
	}

	validRoles := cfg.ValidRoleKeys() // NEW

	issues := struct {
		roleUnknown   []string
		roleNotObject []string
		unknownKeys   []string
		typeErrors    []string
	}{}

	for rawRole, v := range sp {
		roleKey := strings.ToLower(rawRole)
		if _, ok := validRoles[roleKey]; !ok {
			issues.roleUnknown = append(issues.roleUnknown, rawRole)
			continue
		}

		roleMap, ok := v.(map[string]any)
		if !ok {
			issues.roleNotObject = append(issues.roleNotObject, rawRole)
			continue
		}

		var dst Stats
		report := applyRoleMapToStats(&dst, roleMap)

		prefix := roleKey + "."
		for _, k := range report.UnknownKeys {
			issues.unknownKeys = append(issues.unknownKeys, prefix+k)
		}
		for _, te := range report.TypeErrors {
			issues.typeErrors = append(issues.typeErrors, prefix+te)
		}
	}

	if len(issues.roleUnknown)+len(issues.roleNotObject)+len(issues.unknownKeys)+len(issues.typeErrors) == 0 {
		return nil
	}

	sort.Strings(issues.roleUnknown)
	sort.Strings(issues.roleNotObject)
	sort.Strings(issues.unknownKeys)
	sort.Strings(issues.typeErrors)

	err := fmt.Errorf("roles config validation issues: unknown_roles=%v, non_object_roles=%v, unknown_keys=%v, type_errors=%v",
		issues.roleUnknown, issues.roleNotObject, issues.unknownKeys, issues.typeErrors,
	)
	// Behavior unchanged: always return err; Strict decides how upstream handles it.
	return err
}
