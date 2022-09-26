package clearingway

import (
	"fmt"

	"github.com/Veraticus/clearingway/internal/fflogs"
)

var UltimateEncounters = &Encounters{
	Encounters: []*Encounter{
		{
			Name:         "The Unending Coil of Bahamut (Ultimate)",
			Ids:          []int{1060, 1047, 1039},
			Difficulty:   "Ultimate",
			DefaultRoles: false,
		},
		{Name: "The Weapon's Refrain (Ultimate)", Ids: []int{1061, 1048, 1042}, Difficulty: "Ultimate", DefaultRoles: false},
		{Name: "The Epic of Alexander (Ultimate)", Ids: []int{1062, 1050}, Difficulty: "Ultimate", DefaultRoles: false},
		{Name: "Dragonsong's Reprise (Ultimate)", Ids: []int{1065}, Difficulty: "Ultimate", DefaultRoles: false},
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

	e.Roles[ClearedRole].ShouldApply = func(opts *ShouldApplyOpts) (bool, string) {
		for _, id := range e.Ids {
			ranking, ok := opts.Rankings.Rankings[id]
			if !ok {
				continue
			}
			cleared := ranking.Cleared()
			if cleared {
				rank := ranking.RanksByTime()[0]
				return true, fmt.Sprintf("Cleared `%v` with `%v` on <t:%v:F> (%v).",
					e.Name,
					rank.Job.Abbreviation,
					rank.StartTime,
					rank.Report.Url(),
				)
			}
		}

		return false, fmt.Sprintf("Has not cleared %v.", e.Name)
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

func (es *Encounters) Add(e *Encounter) {
	for _, existingEncounter := range es.Encounters {
		if e.Name == existingEncounter.Name {
			continue
		}
	}
	es.Encounters = append(es.Encounters, e)
}

func (es *Encounters) BestDPSRank(rankings *fflogs.Rankings) (*Encounter, *fflogs.Rank) {
	var bestRank *fflogs.Rank
	var bestEncounter *Encounter
	for _, encounter := range es.Encounters {
		for _, encounterId := range encounter.Ids {
			encounterRanking, ok := rankings.Rankings[encounterId]
			if !ok {
				continue
			}
			if !encounterRanking.Cleared() {
				continue
			}

			rank := encounterRanking.BestDPSRank()
			if bestRank == nil || (rank.DPSPercent > bestRank.DPSPercent) {
				bestRank = rank
				bestEncounter = encounter
			}
		}
	}

	return bestEncounter, bestRank
}

func (es *Encounters) WorstDPSRank(rankings *fflogs.Rankings) (*Encounter, *fflogs.Rank) {
	var worstRank *fflogs.Rank
	var worstEncounter *Encounter
	for _, encounter := range es.Encounters {
		for _, encounterId := range encounter.Ids {
			encounterRanking, ok := rankings.Rankings[encounterId]
			if !ok {
				continue
			}
			if !encounterRanking.Cleared() {
				continue
			}

			rank := encounterRanking.WorstDPSRank()
			if worstRank == nil || (rank.DPSPercent < worstRank.DPSPercent) {
				worstRank = rank
				worstEncounter = encounter
			}
		}
	}

	return worstEncounter, worstRank
}

func (es *Encounters) Clears(rankings *fflogs.Rankings) *Encounters {
	encounters := &Encounters{Encounters: []*Encounter{}}
	for _, encounter := range es.Encounters {
		for _, encounterId := range encounter.Ids {
			rankings, ok := rankings.Rankings[encounterId]
			if !ok {
				continue
			}

			if rankings.Cleared() {
				encounters.Encounters = append(encounters.Encounters, encounter)
				continue
			}
		}
	}

	return encounters
}

func (es *Encounters) Names() []string {
	names := []string{}
	for _, e := range es.Encounters {
		names = append(names, e.Name)
	}
	return names
}
