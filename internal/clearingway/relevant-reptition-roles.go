package clearingway

import (
	"fmt"
)

func RelevantRepetitionRoles(encs *Encounters) *Roles {
	roles := &Roles{Roles: []*Role{}}

	for _, enc := range encs.Encounters {
		localEnc := enc
		roles.Roles = append(roles.Roles, &Role{
			Name: "Limbo", Color: 0x808080, Uncomfy: true,
			Description: "Cleared " + localEnc.Name + "... but only once.",
			ShouldApply: func(opts *ShouldApplyOpts) (bool, string) {
				for _, encounter := range opts.Encounters.Encounters {
					if encounter.Name != localEnc.Name {
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
								"Cleared " + localEnc.Name + "... but only once.\nUse `/uncomfy` if you don't want this role.",
							)
					}
				}

				return false, "Cleared " + localEnc.Name + " more than once."
			},
		})

		if localEnc.TotalWeaponsAvailable == 0 {
			continue
		}

		roles.Roles = append(roles.Roles, &Role{
			Name: "Complete", Color: 0xffde00,
			Description: "Cleared " + localEnc.Name + " enough times to have every single weapon.",
			ShouldApply: func(opts *ShouldApplyOpts) (bool, string) {
				for _, encounter := range opts.Encounters.Encounters {
					if encounter.Name != localEnc.Name {
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

					if clears >= localEnc.TotalWeaponsAvailable {
						return true,
							fmt.Sprintf(
								"Cleared " + localEnc.Name + " enough times to have every single weapon.",
							)
					}
				}

				return false, "Has not cleared " + localEnc.Name + " enough times to have every single weapon."
			},
		})
	}

	return roles
}
