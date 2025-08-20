package units

import (
	"fmt"
)

// BuildUnit create a unit by choosing a primary role to override Stats.
func BuildUnit(name string, cost int, traits []string, roles []string, cfg RolesLoader, statOpts ...Option) (Unit, error) {
	if len(roles) == 0 {
		return Unit{}, fmt.Errorf("no role provided for unit %q", name)
	}
	primary := roles[0]

	stats, err := StatsForRole(primary, cfg, statOpts...)
	if err != nil {
		return Unit{}, err
	}

	u := Unit{
		ID:     NewUUID(),
		Name:   name,
		Cost:   cost,
		Traits: traits,
		Roles:  roles,
		Stats:  stats,
	}
	return u, nil
}
