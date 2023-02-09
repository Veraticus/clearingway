package clearingway

import (
	"fmt"
)

func ProgRoles(rs []*ConfigRole, e *Encounter) *Roles {
	roles := &Roles{Roles: []*Role{}}
	for i, r := range rs {
		phase := i
		role := &Role{
			Name: r.Name, Color: r.Color, Type: ProgRole,
			Description: fmt.Sprintf("Reached phase %d (%s) in prog.", phase+1, r.Name),
		}
		roles.Roles = append(roles.Roles, role)
	}

	roles.ShouldApply = func(opts *ShouldApplyOpts) (bool, string, []*Role, []*Role) {
		if opts.Fights == nil || len(opts.Fights.Fights) == 0 {
			return false, "No valid fights found in provided report!", nil, nil
		}

		// Determine if any existing prog role exists
		var existingProgRole *Role
		var existingProgRoleIndex int
		for _, existingRole := range opts.ExistingRoles.Roles {
			if existingRole.Type == ClearedRole {
				existingProgRoleIndex = len(roles.Roles) + 1
				existingProgRole = e.Roles[ClearedRole]
			}

			ok, i := roles.IndexOfRole(existingRole)
			if ok {
				existingProgRole = existingRole
				existingProgRoleIndex = i
			}
		}

		// Find the furthest prog provided in the fights from the report
		furthestFight := opts.Fights.FurthestFight()
		var furthestProgRole *Role
		var furthestProgRoleIndex int

		// Does this report contain a kill?
		if furthestFight.Kill {
			furthestProgRoleIndex = len(roles.Roles) + 1
			furthestProgRole = e.Roles[ClearedRole]
		} else {
			furthestProgRoleIndex = furthestFight.LastPhaseIndex
			furthestProgRole = roles.Roles[furthestProgRoleIndex]
		}

		// Bail out if the furthest prog point in the fight is less than one
		// the user already possesses
		if existingProgRole != nil && furthestProgRoleIndex <= existingProgRoleIndex {
			return false, fmt.Sprintf(
				"You already have a prog role further than the furthest prog in this report! Your existing prog point is `%s` (%s), and the furthest prog point seen by you in this report is `%s` (%s) ⮕ %s",
				existingProgRole.Name,
				existingProgRole.Phase(existingProgRoleIndex+1),
				furthestProgRole.Name,
				furthestProgRole.Phase(furthestProgRoleIndex+1),
				furthestFight.ReportURL(),
			), nil, nil
		}

		// If this fight has the same furthest prog point the user already has,
		// we are done.
		if existingProgRole != nil && existingProgRoleIndex == furthestProgRoleIndex {
			return false, fmt.Sprintf(
				"Your furthest prog point `%s` (%s) is the same as the furthest prog point in this report ⮕ %s",
				existingProgRole.Name,
				existingProgRole.Phase(existingProgRoleIndex+1),
				furthestFight.ReportURL(),
			), nil, nil
		}

		// Looks like we have some real prog to give!
		var lowerRoles []*Role
		if furthestFight.Kill {
			lowerRoles = roles.Roles
		} else {
			lowerRoles = roles.Roles[0:furthestProgRoleIndex]
		}

		return true, fmt.Sprintf(
				"Your furthest prog is now `%s` (%s) ⮕ %s",
				furthestProgRole.Name,
				furthestProgRole.Phase(furthestProgRoleIndex+1),
				furthestFight.ReportURL(),
			), []*Role{
				furthestProgRole,
			}, lowerRoles
	}

	return roles
}
