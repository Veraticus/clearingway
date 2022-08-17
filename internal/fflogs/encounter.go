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
	TotalKills int     `json:"totalKills"`
	Ranks      []*Rank `json:"ranks"`
}

// An encounter could have multiple IDs, since fflogs considers ultimates in
// different expansion to be different encounters.
type Encounters struct {
	Encounters []*Encounter
}

type Encounter struct {
	Name string
	IDs  []int
}

var UltimateEncounters = []*Encounter{
	{Name: "DSR", IDs: []int{1065}},
	{Name: "UCOB", IDs: []int{1060, 1047, 1039}},
	{Name: "UWU", IDs: []int{1061, 1048, 1042}},
	{Name: "TEA", IDs: []int{1062, 1050}},
}

func (e *Encounters) IDs() []int {
	ids := []int{}
	for _, encounter := range e.Encounters {
		ids = append(ids, encounter.IDs...)
	}
	return ids
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

func (r *Rank) Color()
