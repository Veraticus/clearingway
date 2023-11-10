package clearingway

import (
	"fmt"
)

func UltimateRepetitionRoles() *Roles {
	roles := &Roles{Roles: []*Role{}}

	for _, ult := range UltimateEncounters.Encounters {
		roles.Roles = append(roles.Roles, []*Role{
			{
				Name:        ult.The + " Limbo",
				Color:       0x808080,
				Uncomfy:     true,
				Type:        LimboRole,
				Encounter:   ult,
				Description: "Cleared " + ult.Name + "... but only once.",
				ShouldApply: func(opts *ShouldApplyOpts) (bool, string) {
					for _, encounter := range opts.Encounters.Encounters {
						if encounter.Name != ult.Name {
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
									"Cleared " + ult.Name + "... but only once.\nUse `/uncomfy` if you don't want this role.",
								)
						}
					}

					return false, "Cleared " + ult.Name + " more than once."
				},
			},
			{
				Name:        ult.The,
				Color:       0xffde00,
				Type:        CompleteRole,
				Encounter:   ult,
				Description: "Cleared " + ult.Name + " enough times to have every single weapon.",
				ShouldApply: func(opts *ShouldApplyOpts) (bool, string) {
					for _, encounter := range opts.Encounters.Encounters {
						if encounter.Name != ult.Name {
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

						if clears >= ult.TotalWeaponsAvailable {
							return true,
								fmt.Sprintf(
									"Cleared " + ult.Name + " enough times to have every single weapon.",
								)
						}
					}

					return false, "Has not cleared " + ult.Name + " enough times to have every single weapon."
				},
			},
		}...)
	}

	return roles
}
