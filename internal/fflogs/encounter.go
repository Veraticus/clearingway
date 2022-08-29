package fflogs

import (
	"sort"
)

// An EncounterRankings contains many encounter rankings,
// indexed by the encounter ID of the fight.
type EncounterRankings struct {
	Encounters map[int]*EncounterRanking
}

type EncounterRanking struct {
	Error      string  `json:"error"`
	TotalKills int     `json:"totalKills"`
	Ranks      []*Rank `json:"ranks"`
}

type Rank struct {
	Percent float64 `json:"rankPercent"`
}

// An encounter could have multiple IDs, since fflogs considers ultimates in
// different expansion to be different encounters.
type Encounters struct {
	Encounters []*Encounter `yaml:"encounters"`
}

type Encounter struct {
	Name         string `yaml:"name"`
	Difficulty   string `yaml:"difficulty"`
	IDs          []int  `yaml:"ids"`
	CreateRoles  bool   `yaml:"createRoles"`
	ClearedColor int    `yaml:"clearedColor"`
}

var UltimateEncounters = &Encounters{
	Encounters: []*Encounter{
		{Name: "DSR", IDs: []int{1065}, Difficulty: "Ultimate", CreateRoles: false},
		{Name: "UCOB", IDs: []int{1060, 1047, 1039}, Difficulty: "Ultimate", CreateRoles: false},
		{Name: "UWU", IDs: []int{1061, 1048, 1042}, Difficulty: "Ultimate", CreateRoles: false},
		{Name: "TEA", IDs: []int{1062, 1050}, Difficulty: "Ultimate", CreateRoles: false},
	},
}

func (e *Encounters) BestRankForEncounterRankings(ers *EncounterRankings) *Rank {
	var bestRank *Rank
	for _, encounter := range e.Encounters {
		for _, encounterId := range encounter.IDs {
			encounterRanking, ok := ers.Encounters[encounterId]
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

func (e *Encounters) TotalClearsFromEncounterRankings(ers *EncounterRankings) int {
	clears := map[string]bool{}
	for _, encounter := range e.Encounters {
		for _, encounterId := range encounter.IDs {
			encounterRanking, ok := ers.Encounters[encounterId]
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

func (er *EncounterRanking) Cleared() bool {
	return er.TotalKills > 0
}

func (er *EncounterRanking) BestRank() *Rank {
	ranks := make([]*Rank, len(er.Ranks))
	copy(ranks, er.Ranks)
	sort.SliceStable(ranks, func(i, j int) bool { return ranks[i].Percent > ranks[j].Percent })
	return ranks[0]
}

func (e *Encounter) DifficultyInt() int {
	if e.Difficulty == "Savage" {
		return 101
	}

	return 100
}
