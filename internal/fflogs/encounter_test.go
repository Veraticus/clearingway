package fflogs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBestRank(t *testing.T) {
	bestRank := &Rank{RankPercent: 90.1384}
	encounterRankings := &EncounterRankings{
		Ranks: []*Rank{
			bestRank,
			{RankPercent: 10.13},
			{RankPercent: 50.154},
		},
	}
	assert.Equal(t, bestRank, encounterRankings.BestRank())
}
