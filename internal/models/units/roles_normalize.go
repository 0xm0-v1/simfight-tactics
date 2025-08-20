package units

import "strings"

// detectRoleKey parses a free-form label into (roleKey, damageType, ok).
// It now receives allowed tokens from the caller (data-driven).
func detectRoleKey(raw string, validRoles, damageTokens map[string]struct{}) (string, string, bool) {
	s := strings.ToLower(strings.TrimSpace(raw))
	if s == "" {
		return "", "", false
	}
	tokens := strings.FieldsFunc(s, func(r rune) bool {
		switch r {
		case ' ', '\t', '-', '_', '/':
			return true
		default:
			return false
		}
	})

	var damageType string
	var roleKey string

	for _, tok := range tokens {
		if _, ok := damageTokens[tok]; ok {
			if damageType != "" {
				return "", "", false // multiple damage tokens → invalid
			}
			damageType = tok
			continue
		}
		if _, ok := validRoles[tok]; ok {
			roleKey = tok
			break
		}
		// Unknown token before role ⇒ reject (keeps previous semantics)
		return "", "", false
	}
	if roleKey == "" {
		return "", "", false
	}
	return roleKey, damageType, true
}
