package fflogs

import (
	"sort"
)

type EncounterRankings struct {
	TotalKills int     `json:"totalKills"`
	Ranks      []*Rank `json:"ranks"`
}

type Rank struct {
	RankPercent float64 `json:"rankPercent"`
}

func (er *EncounterRankings) Cleared() bool {
	return er.TotalKills > 0
}

func (er *EncounterRankings) BestRank() *Rank {
	ranks := make([]*Rank, len(er.Ranks))
	copy(ranks, er.Ranks)
	sort.SliceStable(ranks, func(i, j int) bool { return ranks[i].RankPercent > ranks[j].RankPercent })
	return ranks[0]
}
