package clearingway

import (
	"fmt"

	"github.com/Veraticus/clearingway/internal/ffxiv"
)

type Datacenters struct {
	Datacenters map[string]*Datacenter
}

type Datacenter struct {
	Datacenter string
	Role       *Role
}

func (ds *Datacenters) Init(c []*ConfigDatacenter) {
	ds.Datacenters = map[string]*Datacenter{}
	if len(c) == 0 {
		ds.Datacenters["Aether"] = &Datacenter{Datacenter: "Aether", Role: &Role{Name: "Aether", Color: 0x71368a}}
		ds.Datacenters["Primal"] = &Datacenter{Datacenter: "Primal", Role: &Role{Name: "Primal", Color: 0x992d22}}
		ds.Datacenters["Crystal"] = &Datacenter{Datacenter: "Crystal", Role: &Role{Name: "Crystal", Color: 0x206694}}
	} else {
		for _, d := range c {
			ds.Datacenters[d.Datacenter] = &Datacenter{Datacenter: d.Datacenter, Role: &Role{Name: d.Name, Color: d.Color}}
		}
	}

	for _, d := range ds.Datacenters {
		worlds, err := ffxiv.WorldsForDatacenter(d.Datacenter)
		if err != nil {
			panic(fmt.Sprintf("Bad datacenter specified: %+v!", d.Datacenter))
		}
		name := d.Datacenter

		d.Role.ShouldApply = func(opts *ShouldApplyOpts) (bool, string) {
			_, ok := worlds[opts.Character.World]
			if ok {
				return true, fmt.Sprintf(
					"Character `%v (%v)` is in %v.",
					opts.Character.Name(),
					opts.Character.World,
					name,
				)
			}
			return false, fmt.Sprintf(
				"Character `%v (%v)` is not in %v.",
				opts.Character.Name(),
				opts.Character.World,
				name,
			)
		}
	}
}

func (ds *Datacenters) AllRoles() *Roles {
	roles := &Roles{Roles: []*Role{}}

	for _, datacenter := range ds.Datacenters {
		roles.Roles = append(roles.Roles, datacenter.Role)
	}

	return roles
}
