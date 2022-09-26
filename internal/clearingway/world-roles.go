package clearingway

import (
	"fmt"

	"github.com/Veraticus/clearingway/internal/ffxiv"
)

func WorldRoles() *Roles {
	return &Roles{Roles: []*Role{
		{
			Name: "Aether", Color: 0x71368a,
			ShouldApply: func(opts *ShouldApplyOpts) (bool, string) {
				_, ok := ffxiv.AetherWorlds[opts.Character.World]
				if ok {
					return true, fmt.Sprintf("Charater `%v (%v)` is in Aether.", opts.Character.Name(), opts.Character.World)
				}
				return false, fmt.Sprintf("Character `%v (%v)` is not in Aether.", opts.Character.Name(), opts.Character.World)
			},
		},
		{
			Name: "Crystal", Color: 0x206694,
			ShouldApply: func(opts *ShouldApplyOpts) (bool, string) {
				_, ok := ffxiv.CrystalWorlds[opts.Character.World]
				if ok {
					return true, fmt.Sprintf("Charater `%v (%v)` is in Crystal.", opts.Character.Name(), opts.Character.World)
				}
				return false, fmt.Sprintf("Character `%v (%v)` is not in Crystal.", opts.Character.Name(), opts.Character.World)
			},
		},
		{
			Name: "Primal", Color: 0x992d22,
			ShouldApply: func(opts *ShouldApplyOpts) (bool, string) {
				_, ok := ffxiv.PrimalWorlds[opts.Character.World]
				if ok {
					return true, fmt.Sprintf("Charater `%v (%v)` is in Primal.", opts.Character.Name(), opts.Character.World)
				}
				return false, fmt.Sprintf("Character `%v (%v)` is not in Primal.", opts.Character.Name(), opts.Character.World)
			},
		},
	}}
}
