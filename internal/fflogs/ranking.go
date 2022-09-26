package fflogs

import (
	"fmt"
	"sort"

	"github.com/Veraticus/clearingway/internal/ffxiv"
)

// Rankings contains many rankings for different encounters,
// indexed by the encounter ID of the fight.
type Rankings struct {
	Rankings map[int]*Ranking
}

type Metric string

const (
	Dps Metric = "dps"
	Hps Metric = "hps"
)

type Ranking struct {
	Error      string  `json:"error"`
	TotalKills int     `json:"totalKills"`
	Ranks      []*Rank `json:"ranks"`
}

type Rank struct {
	Percent   float64 `json:"rankPercent"`
	Spec      string  `json:"spec"`
	StartTime int     `json:"startTime"`
	Report    Report  `json:"report"`
	Metric    Metric
	Job       *ffxiv.Job
}

type Report struct {
	Code    string `json:"code"`
	FightId int    `json:"fightID"`
}

func (r *Ranking) Cleared() bool {
	return r.TotalKills > 0
}

func (r *Ranking) DPSRanksByPercent() []*Rank {
	ranks := []*Rank{}
	for _, r := range r.Ranks {
		if r.Metric == Dps {
			ranks = append(ranks, r)
		}
	}
	sort.SliceStable(ranks, func(i, j int) bool { return ranks[i].Percent > ranks[j].Percent })
	return ranks
}

func (r *Rank) BestParseString(encounterName string) string {
	return fmt.Sprintf(
		"Best parse was *%v* with `%v` in `%v` on <t:%v:F> (%v).",
		r.Percent,
		r.Job.Abbreviation,
		encounterName,
		r.StartTime,
		r.Report.Url(),
	)
}

func (r *Ranking) BestDPSRank() *Rank {
	return r.DPSRanksByPercent()[0]
}

func (r *Ranking) WorstDPSRank() *Rank {
	sortedRanks := r.DPSRanksByPercent()
	return sortedRanks[len(sortedRanks)-1]
}

func (r *Ranking) HPSRanksByPercent() []*Rank {
	ranks := []*Rank{}
	for _, r := range r.Ranks {
		if r.Metric == Hps {
			ranks = append(ranks, r)
		}
	}
	sort.SliceStable(ranks, func(i, j int) bool { return ranks[i].Percent > ranks[j].Percent })
	return ranks
}

func (r *Ranking) BestHPSRank() *Rank {
	return r.HPSRanksByPercent()[0]
}

func (r *Ranking) WorstHPSRank() *Rank {
	sortedRanks := r.HPSRanksByPercent()
	return sortedRanks[len(sortedRanks)-1]
}

func (r *Ranking) RanksByTime() []*Rank {
	ranks := make([]*Rank, len(r.Ranks))
	copy(ranks, r.Ranks)
	sort.SliceStable(ranks, func(i, j int) bool { return ranks[i].StartTime < ranks[j].StartTime })
	return ranks
}

func (r *Report) Url() string {
	return fmt.Sprintf("https://www.fflogs.com/reports/%v#fight=%v", r.Code, r.FightId)
}
