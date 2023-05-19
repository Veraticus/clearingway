package clearingway

import (
	"fmt"
)

func UltimateReptitionRoles() *Roles {
	roles := &Roles{Roles: []*Role{}}

	for _, ult := range UltimateEncounters.Encounters {
		localUlt := ult
		roles.Roles = append(roles.Roles, []*Role{
			{
				Name: localUlt.The + " Limbo", Color: 0x808080, Uncomfy: true,
				Description: "Cleared " + localUlt.Name + "... but only once.",
				ShouldApply: func(opts *ShouldApplyOpts) (bool, string) {
					for _, encounter := range opts.Encounters.Encounters {
						if encounter.Name != localUlt.Name {
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
									"Cleared " + localUlt.Name + "... but only once.\nUse `/uncomfy` if you don't want this role.",
								)
						}
					}

					return false, "Cleared " + localUlt.Name + " more than once."
				},
			},
			{
				Name: localUlt.The, Color: 0xffde00,
				Description: "Cleared " + localUlt.Name + " enough times to have every single weapon.",
				ShouldApply: func(opts *ShouldApplyOpts) (bool, string) {
					for _, encounter := range opts.Encounters.Encounters {
						if encounter.Name != localUlt.Name {
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

						if clears >= localUlt.TotalWeaponsAvailable {
							return true,
								fmt.Sprintf(
									"Cleared " + localUlt.Name + " enough times to have every single weapon.",
								)
						}
					}

					return false, "Has not cleared " + localUlt.Name + " enough times to have every single weapon."
				},
			},
		}...)
	}

	return roles
}
