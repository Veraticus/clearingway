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
	Dps Metric = "rdps"
	Hps Metric = "hps"
)

type Ranking struct {
	Error      string  `json:"error"`
	TotalKills int     `json:"totalKills"`
	Metric     Metric  `json:"metric"`
	Ranks      []*Rank `json:"ranks"`
}

type Rank struct {
	RankPercent float64 `json:"rankPercent"`
	Spec        string  `json:"spec"`
	StartTime   int     `json:"startTime"`
	Report      Report  `json:"report"`
	Job         *ffxiv.Job

	DPSParseFound bool
	HPSParseFound bool

	DPSPercent float64
	HPSPercent float64
}

type Report struct {
	Code      string `json:"code"`
	StartTime int    `json:"startTime"`
	FightId   int    `json:"fightID"`
}

func (rs *Rankings) Add(id int, r *Ranking) error {
	existingRankings, ok := rs.Rankings[id]

	if !ok {
		for _, rank := range r.Ranks {
			if r.Metric == Hps {
				rank.HPSPercent = rank.RankPercent
				rank.HPSParseFound = true
			} else if r.Metric == Dps {
				rank.DPSPercent = rank.RankPercent
				rank.DPSParseFound = true
			}
			j, ok := ffxiv.Jobs[rank.Spec]
			if !ok {
				return fmt.Errorf("Could not find job %s", rank.Spec)
			}
			rank.Job = j
		}
		rs.Rankings[id] = r
		return nil
	}

	if r.TotalKills != 0 {
		rs.Rankings[id].TotalKills = r.TotalKills
	}

	for _, newRank := range r.Ranks {
		found := false
		for _, existingRank := range existingRankings.Ranks {
			if existingRank.SameFight(newRank) {
				found = true
				if r.Metric == Hps {
					existingRank.HPSPercent = newRank.RankPercent
					existingRank.HPSParseFound = true
				} else if r.Metric == Dps {
					existingRank.DPSPercent = newRank.RankPercent
					existingRank.DPSParseFound = true
				}
				continue
			}
		}
		if !found {
			if r.Metric == Hps {
				newRank.HPSPercent = newRank.RankPercent
				newRank.HPSParseFound = true
			} else if r.Metric == Dps {
				newRank.DPSPercent = newRank.RankPercent
				newRank.DPSParseFound = true
			}
			j, ok := ffxiv.Jobs[newRank.Spec]
			if !ok {
				return fmt.Errorf("Could not find job %s", newRank.Spec)
			}
			newRank.Job = j
			rs.Rankings[id].Ranks = append(rs.Rankings[id].Ranks, newRank)
		}
	}

	return nil
}

func (r *Rank) SameFight(o *Rank) bool {
	return r.StartTime == o.StartTime
}

func (r *Ranking) Cleared() bool {
	return r.TotalKills > 0
}

func (r *Ranking) RanksByDPSPercent() []*Rank {
	ranks := []*Rank{}
	for _, r := range r.Ranks {
		if r.DPSParseFound {
			ranks = append(ranks, r)
		}
	}
	sort.SliceStable(ranks, func(i, j int) bool { return ranks[i].DPSPercent > ranks[j].DPSPercent })
	return ranks
}

func (r *Ranking) BestDPSRank() *Rank {
	return r.RanksByDPSPercent()[0]
}

func (r *Ranking) WorstDPSRank() *Rank {
	sortedRanks := r.RanksByDPSPercent()
	if len(sortedRanks) == 0 {
		return nil
	}
	return sortedRanks[len(sortedRanks)-1]
}

func (r *Ranking) RanksByHPSPercent() []*Rank {
	ranks := []*Rank{}
	for _, r := range r.Ranks {
		if r.HPSParseFound {
			ranks = append(ranks, r)
		}
	}
	sort.SliceStable(ranks, func(i, j int) bool { return ranks[i].HPSPercent > ranks[j].HPSPercent })
	return ranks
}

func (r *Ranking) BestHPSRank() *Rank {
	return r.RanksByHPSPercent()[0]
}

func (r *Ranking) WorstHPSRank() *Rank {
	sortedRanks := r.RanksByHPSPercent()
	if len(sortedRanks) == 0 {
		return nil
	}
	return sortedRanks[len(sortedRanks)-1]
}

func (r *Ranking) RanksByTime() []*Rank {
	ranks := make([]*Rank, len(r.Ranks))
	copy(ranks, r.Ranks)
	sort.SliceStable(ranks, func(i, j int) bool { return ranks[i].StartTime < ranks[j].StartTime })
	return ranks
}

func (r *Rank) BestDPSParseString(encounterName string) string {
	return fmt.Sprintf(
		"Best parse was **%v** with `%v` in `%v` on <t:%v:F> (%v).",
		r.DPSPercentString(),
		r.Job.Abbreviation,
		encounterName,
		r.UnixTime(),
		r.Report.Url(),
	)
}

func (r *Rank) UnixTime() int {
	return r.StartTime / 1000
}

func (r *Rank) DPSPercentString() string {
	return fmt.Sprintf("%.2f", r.DPSPercent)
}

func (r *Rank) HPSPercentString() string {
	return fmt.Sprintf("%.2f", r.HPSPercent)
}

func (r *Report) Url() string {
	return fmt.Sprintf("https://www.fflogs.com/reports/%v#fight=%v", r.Code, r.FightId)
}
