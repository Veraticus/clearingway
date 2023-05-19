package clearingway

import (
	"fmt"
)

func RelevantReptitionRoles(encs *Encounters) *Roles {
	roles := &Roles{Roles: []*Role{}}

	for _, enc := range encs.Encounters {
		roles.Roles = append(roles.Roles, &Role{
			Name: "Limbo", Color: 0x808080, Uncomfy: true,
			Description: "Cleared " + enc.Name + "... but only once.",
			ShouldApply: func(opts *ShouldApplyOpts) (bool, string) {
				for _, encounter := range opts.Encounters.Encounters {
					if encounter.Name != enc.Name {
						continue
					}
					clears := 0

					for _, encounterId := range encounter.Ids {
						ranking, ok := opts.Rankings.Rankings[encounterId]

						if !ok {
							continue
						}
						if !ranking.Cleared() {
							continue
						}

						clears = clears + ranking.TotalKills
					}

					if clears == 1 {
						return true,
							fmt.Sprintf(
								"Cleared " + enc.Name + "... but only once.\nUse `/uncomfy` if you don't want this role.",
							)
					}
				}

				return false, "Cleared " + enc.Name + " more than once."
			},
		})

		if enc.TotalWeaponsAvailable == 0 {
			continue
		}

		roles.Roles = append(roles.Roles, &Role{
			Name: "Complete", Color: 0xffde00,
			Description: "Cleared " + enc.Name + " enough times to have every single weapon.",
			ShouldApply: func(opts *ShouldApplyOpts) (bool, string) {
				for _, encounter := range opts.Encounters.Encounters {
					if encounter.Name != enc.Name {
						continue
					}
					clears := 0

					for _, encounterId := range encounter.Ids {
						ranking, ok := opts.Rankings.Rankings[encounterId]

						if !ok {
							continue
						}
						if !ranking.Cleared() {
							continue
						}

						clears = clears + ranking.TotalKills
					}

					if clears >= enc.TotalWeaponsAvailable {
						return true,
							fmt.Sprintf(
								"Cleared " + enc.Name + " enough times to have every single weapon.",
							)
					}
				}

				return false, "Has not cleared " + enc.Name + " enough times to have every single weapon."
			},
		})
	}

	return roles
}
