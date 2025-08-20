package units

import "testing"

func TestDetectRoleKey_TokensAndDelimiters(t *testing.T) {
	t.Parallel()
	valid := map[string]struct{}{"tank": {}, "assassin": {}, "fighter": {}}
	dmg := map[string]struct{}{"attack": {}, "magic": {}, "hybrid": {}}

	type tc struct {
		in, role, dmg string
		ok            bool
	}
	cases := []tc{
		{"Attack Tank", "tank", "attack", true},
		{" MAGIC_assassin ", "assassin", "magic", true},
		{"hybrid--fighter", "fighter", "hybrid", true},
		{"attack magic tank", "", "", false}, // multiple damage tokens
		{"unknown tank", "", "", false},      // unknown token before role
		{"tank", "tank", "", true},           // damage optional
	}
	for _, c := range cases {
		r, d, ok := detectRoleKey(c.in, valid, dmg)
		if ok != c.ok || r != c.role || d != c.dmg {
			t.Fatalf("%q → got (r=%q,d=%q,ok=%v) want (r=%q,d=%q,ok=%v)",
				c.in, r, d, ok, c.role, c.dmg, c.ok)
		}
	}
}
