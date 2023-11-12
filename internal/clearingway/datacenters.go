package clearingway

import (
	"fmt"

	"github.com/Veraticus/clearingway/internal/ffxiv"
)

type PhysicalDatacenters struct {
	PhysicalDatacenters map[string]*PhysicalDatacenter
}

type PhysicalDatacenter struct {
	Name               string
	LogicalDatacenters map[string]*LogicalDatacenter
}

type LogicalDatacenter struct {
	Name string
	Role *Role
}

func (pds *PhysicalDatacenters) Init(c []*ConfigPhysicalDatacenter) {
	for _, cpd := range c {
		ffxivPhysicalDatacenter := ffxiv.PhysicalDatacenterForAbbreviation(cpd.Name)
		ffxivPhysicalDatacenterName := ffxivPhysicalDatacenter.Name
		pd := &PhysicalDatacenter{
			Name:               ffxivPhysicalDatacenter.Name,
			LogicalDatacenters: map[string]*LogicalDatacenter{},
		}

		for _, ffxivLogicalDatacenter := range ffxivPhysicalDatacenter.LogicalDatacenters {
			ffxivLogicalDatacenterName := ffxivLogicalDatacenter.Name
			ffxivLogicalDatacenterWorlds := ffxivLogicalDatacenter.Worlds
			ld := &LogicalDatacenter{
				Name: ffxivLogicalDatacenterName,
				Role: &Role{
					Name: ffxivLogicalDatacenterName,
					Description: fmt.Sprintf(
						"Is in the %v datacenter in %v.",
						ffxivLogicalDatacenterName,
						ffxivPhysicalDatacenterName,
					),
					Color: ffxivLogicalDatacenter.Color,
					ShouldApply: func(opts *ShouldApplyOpts) (bool, string) {
						_, ok := ffxivLogicalDatacenterWorlds[opts.Character.World]
						if ok {
							return true, fmt.Sprintf(
								"Character `%v (%v)` is in the %v datacenter in %v.",
								opts.Character.Name(),
								opts.Character.World,
								ffxivLogicalDatacenterName,
								ffxivPhysicalDatacenterName,
							)
						}
						return false, fmt.Sprintf(
							"Character `%v (%v)` is not in the %v datacenter in %v.",
							opts.Character.Name(),
							opts.Character.World,
							ffxivLogicalDatacenterName,
							ffxivPhysicalDatacenterName,
						)
					},
				},
			}
			pd.LogicalDatacenters[ffxivLogicalDatacenterName] = ld
		}

		if len(cpd.LogicalDatacenters) != 0 {
			for _, cld := range cpd.LogicalDatacenters {
				from := pd.LogicalDatacenters[cld.From]
				if cld.To != "" {
					from.Role.Name = cld.To
				}
				if cld.Color != 0 {
					from.Role.Color = cld.Color
				}
				if cld.Hoist {
					from.Role.Hoist = true
				}
			}
		}

		pds.PhysicalDatacenters[cpd.Name] = pd
	}
}

func (pds *PhysicalDatacenters) AllRoles() *Roles {
	roles := &Roles{Roles: []*Role{}}

	for _, physicalDatacenter := range pds.PhysicalDatacenters {
		for _, logicalDatacenter := range physicalDatacenter.LogicalDatacenters {
			roles.Roles = append(roles.Roles, logicalDatacenter.Role)
		}
	}

	return roles
}
