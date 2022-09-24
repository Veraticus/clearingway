package clearingway

import (
	"github.com/Veraticus/clearingway/internal/fflogs"
)

var UltimateEncounters = &Encounters{
	Encounters: []*Encounter{
		{Name: "DSR", Ids: []int{1065}, Difficulty: "Ultimate", DefaultRoles: false},
		{Name: "UCOB", Ids: []int{1060, 1047, 1039}, Difficulty: "Ultimate", DefaultRoles: false},
		{Name: "UWU", Ids: []int{1061, 1048, 1042}, Difficulty: "Ultimate", DefaultRoles: false},
		{Name: "TEA", Ids: []int{1062, 1050}, Difficulty: "Ultimate", DefaultRoles: false},
	},
}

type Encounters struct {
	Encounters []*Encounter
}

type Encounter struct {
	Name         string `yaml:"name"`
	Difficulty   string `yaml:"difficulty"`
	DefaultRoles bool   `yaml:"defaultRoles"`
	Ids          []int  `yaml:"ids"`
	Roles        map[RoleType]*Role
}

func (e *Encounter) Init(c *ConfigEncounter) {
	e.Ids = c.Ids
	e.Name = c.Name
	e.Difficulty = c.Difficulty
	e.DefaultRoles = c.DefaultRoles
	e.Roles = map[RoleType]*Role{}

	if e.DefaultRoles {
		e.Roles[PfRole] = &Role{Name: e.Name + "-PF", Color: 0x11806a, Type: PfRole}
		e.Roles[ReclearRole] = &Role{Name: e.Name + "-Reclear", Color: 0x11806a, Type: ReclearRole}
		e.Roles[ParseRole] = &Role{Name: e.Name + "-Parse", Color: 0x11806a, Type: ParseRole}
		e.Roles[ClearedRole] = &Role{Name: e.Name + "-Cleared", Color: 0x11806a, Type: ClearedRole}
	}

	for _, configRole := range c.ConfigRoles {
		roleType := RoleType(configRole.Type)
		role, ok := e.Roles[roleType]
		if !ok {
			role = &Role{Type: roleType}
			e.Roles[roleType] = role
		}
		if len(configRole.Name) != 0 {
			role.Name = configRole.Name
		}
		if configRole.Color != 0 {
			role.Color = configRole.Color
		}
	}

	e.Roles[ClearedRole].ShouldApply = func(opts *ShouldApplyOpts) bool {
		for _, id := range e.Ids {
			encounterRanking, ok := opts.Rankings.Rankings[id]
			if !ok {
				continue
			}
			cleared := encounterRanking.Cleared()
			if cleared {
				return true
			}
		}

		return false
	}
}

func (e *Encounter) DifficultyInt() int {
	if e.Difficulty == "Savage" {
		return 101
	}

	return 100
}

func (es *Encounters) Roles() *Roles {
	roles := &Roles{Roles: []*Role{}}

	for _, encounter := range es.Encounters {
		for _, role := range encounter.Roles {
			roles.Roles = append(roles.Roles, role)
		}
	}

	return roles
}

func (es *Encounters) BestRank(rankings *fflogs.Rankings) *fflogs.Rank {
	var bestRank *fflogs.Rank
	for _, encounter := range es.Encounters {
		for _, encounterId := range encounter.Ids {
			encounterRanking, ok := rankings.Rankings[encounterId]
			if !ok {
				continue
			}
			if !encounterRanking.Cleared() {
				continue
			}

			rank := encounterRanking.BestRank()
			if bestRank == nil || (rank.Percent > bestRank.Percent) {
				bestRank = rank
			}
		}
	}

	return bestRank
}

func (es *Encounters) WorstRank(rankings *fflogs.Rankings) *fflogs.Rank {
	var worstRank *fflogs.Rank
	for _, encounter := range es.Encounters {
		for _, encounterId := range encounter.Ids {
			encounterRanking, ok := rankings.Rankings[encounterId]
			if !ok {
				continue
			}
			if !encounterRanking.Cleared() {
				continue
			}

			rank := encounterRanking.WorstRank()
			if worstRank == nil || (rank.Percent < worstRank.Percent) {
				worstRank = rank
			}
		}
	}

	return worstRank
}

func (es *Encounters) TotalClears(rankings *fflogs.Rankings) int {
	clears := map[string]bool{}
	for _, encounter := range es.Encounters {
		for _, encounterId := range encounter.Ids {
			encounterRanking, ok := rankings.Rankings[encounterId]
			if !ok {
				continue
			}
			if encounterRanking.Cleared() {
				clears[encounter.Name] = true
				continue
			}
		}
	}

	totalClears := 0
	for _, clear := range clears {
		if clear == true {
			totalClears = totalClears + 1
		}
	}

	return totalClears
}
