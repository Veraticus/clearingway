package fflogs

import (
	"sort"
)

// Rankings contains many rankings for different encounters,
// indexed by the encounter ID of the fight.
type Rankings struct {
	Rankings map[int]*Ranking
}

type Ranking struct {
	Error      string  `json:"error"`
	TotalKills int     `json:"totalKills"`
	Ranks      []*Rank `json:"ranks"`
}

type Rank struct {
	Percent float64 `json:"rankPercent"`
}

func (r *Ranking) Cleared() bool {
	return r.TotalKills > 0
}

func (r *Ranking) BestRank() *Rank {
	ranks := make([]*Rank, len(r.Ranks))
	copy(ranks, r.Ranks)
	sort.SliceStable(ranks, func(i, j int) bool { return ranks[i].Percent > ranks[j].Percent })
	return ranks[0]
}
