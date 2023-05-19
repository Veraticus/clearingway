package clearingway

import (
	"fmt"

	"github.com/Veraticus/clearingway/internal/fflogs"
)

var UltimateEncounters = &Encounters{
	Encounters: []*Encounter{
		{
			Name:                  "The Unending Coil of Bahamut (Ultimate)",
			Ids:                   []int{1060, 1047, 1039},
			Difficulty:            "Ultimate",
			DefaultRoles:          false,
			TotalWeaponsAvailable: 15,
			The:                   "Legendary",
		},
		{
			Name:                  "The Weapon's Refrain (Ultimate)",
			Ids:                   []int{1061, 1048, 1042},
			Difficulty:            "Ultimate",
			DefaultRoles:          false,
			TotalWeaponsAvailable: 15,
			The:                   "Ultimate",
		},
		{
			Name:                  "The Epic of Alexander (Ultimate)",
			Ids:                   []int{1062, 1050},
			Difficulty:            "Ultimate",
			DefaultRoles:          false,
			TotalWeaponsAvailable: 17,
			The:                   "Perfect",
		},
		{
			Name:                  "Dragonsong's Reprise (Ultimate)",
			Ids:                   []int{1065},
			Difficulty:            "Ultimate",
			DefaultRoles:          false,
			TotalWeaponsAvailable: 19,
			The:                   "Heavenly",
		},
		{
			Name:                  "The Omega Protocol (Ultimate)",
			Ids:                   []int{1068},
			Difficulty:            "Ultimate",
			DefaultRoles:          false,
			TotalWeaponsAvailable: 19,
			The:                   "Alpha",
		},
	},
}

type Encounters struct {
	Encounters []*Encounter
}

type Encounter struct {
	Name                  string `yaml:"name"`
	Difficulty            string `yaml:"difficulty"`
	DefaultRoles          bool   `yaml:"defaultRoles"`
	Ids                   []int  `yaml:"ids"`
	TotalWeaponsAvailable int
	Roles                 map[RoleType]*Role
	ProgRoles             *Roles
	The                   string
}

func (e *Encounter) Init(c *ConfigEncounter) {
	e.Ids = c.Ids
	e.Name = c.Name
	e.Difficulty = c.Difficulty
	e.DefaultRoles = c.DefaultRoles
	e.TotalWeaponsAvailable = c.TotalWeaponsAvailable
	e.The = c.The
	e.Roles = map[RoleType]*Role{}

	if e.DefaultRoles {
		e.Roles[PfRole] = &Role{
			Name:  e.Name + "-PF",
			Color: 0x11806a,
			Type:  PfRole,
		}
		e.Roles[ReclearRole] = &Role{
			Name:  e.Name + "-Reclear",
			Color: 0x11806a,
			Type:  ReclearRole,
		}
		e.Roles[ParseRole] = &Role{
			Name:  e.Name + "-Parse",
			Color: 0x11806a,
			Type:  ParseRole,
		}
		e.Roles[ClearedRole] = &Role{
			Name:  e.Name + "-Cleared",
			Color: 0x11806a,
			Type:  ClearedRole,
		}
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
		if len(configRole.Description) != 0 {
			role.Description = configRole.Description
		}
		if configRole.Color != 0 {
			role.Color = configRole.Color
		}
		if configRole.Hoist {
			role.Hoist = true
		}
		if configRole.Mention {
			role.Mention = true
		}
	}

	for roleType, role := range e.Roles {
		if len(role.Description) != 0 {
			continue
		}
		switch roleType {
		case PfRole:
			role.Description = "Wants to PF " + e.Name + "."
		case ReclearRole:
			role.Description = "Wants to reclear " + e.Name + "."
		case ParseRole:
			role.Description = "Wants to parse " + e.Name + "."
		case ClearedRole:
			role.Description = "Cleared " + e.Name + "."
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
					rank.UnixTime(),
					rank.Report.Url(),
				)
			}
		}

		return false, fmt.Sprintf("Has not cleared %v.", e.Name)
	}

	if c.ConfigProg != nil {
		e.ProgRoles = ProgRoles(c.ConfigProg, e)
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
		if encounter.ProgRoles != nil {
			roles.Roles = append(roles.Roles, encounter.ProgRoles.Roles...)
		}
	}

	return roles
}

func (es *Encounters) Add(e *Encounter) {
	for _, existingEncounter := range es.Encounters {
		if e.Name == existingEncounter.Name {
			return
		}
	}

	es.Encounters = append(es.Encounters, e)
}

func (es *Encounters) BestDPSRank(rankings *fflogs.Rankings) (*Encounter, *fflogs.Rank) {
	var bestRank *fflogs.Rank
	var bestEncounter *Encounter
	for _, encounter := range es.Encounters {
		for _, encounterId := range encounter.Ids {
			ranking, ok := rankings.Rankings[encounterId]
			if !ok {
				continue
			}
			if !ranking.Cleared() {
				continue
			}

			rank := ranking.BestDPSRank()
			if rank == nil {
				continue
			}
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
			ranking, ok := rankings.Rankings[encounterId]
			if !ok {
				continue
			}
			if !ranking.Cleared() {
				continue
			}

			rank := ranking.WorstDPSRank()
			if rank == nil {
				continue
			}
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
				encounters.Add(encounter)
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

func (es *Encounters) ForName(name string) *Encounter {
	for _, e := range es.Encounters {
		if e.Name == name {
			return e
		}
	}
	return nil
}

func (e *Encounter) Ranks(rankings *fflogs.Rankings) []*fflogs.Rank {
	ranks := []*fflogs.Rank{}
	for id, ranking := range rankings.Rankings {
		for _, encounterId := range e.Ids {
			if encounterId == id {
				ranks = append(ranks, ranking.Ranks...)
			}
		}
	}

	return ranks
}

func (e *Encounter) Fights(fights *fflogs.Fights) []*fflogs.Fight {
	fs := []*fflogs.Fight{}
	for _, fight := range fights.Fights {
		for _, encounterId := range e.Ids {
			if encounterId == fight.EncounterID {
				fs = append(fs, fight)
			}
		}
	}

	return fs
}
