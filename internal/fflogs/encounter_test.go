package fflogs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBestRank(t *testing.T) {
	bestRank := &Rank{Percent: 90.1384}
	encounterRanking := &EncounterRanking{
		TotalKills: 2,
		Ranks: []*Rank{
			bestRank,
			{Percent: 10.13},
			{Percent: 50.154},
		},
	}
	assert.Equal(t, bestRank, encounterRanking.BestRank())
}

func TestBestRankForEncounterRankings_FindBest(t *testing.T) {
	encounter1 := &Encounter{
		Name: "Encounter 1",
		IDs:  []int{1},
	}
	bestRank := &Rank{Percent: 90.1384}
	encounterRanking1 := &EncounterRanking{
		TotalKills: 2,
		Ranks: []*Rank{
			bestRank,
			{Percent: 10.13},
			{Percent: 50.154},
		},
	}

	encounter2 := &Encounter{
		Name: "Encounter 2",
		IDs:  []int{2},
	}
	encounterRanking2 := &EncounterRanking{
		TotalKills: 3,
		Ranks: []*Rank{
			{Percent: 1.1},
			{Percent: 2.2},
		},
	}

	encounterRankings := &EncounterRankings{
		Encounters: map[int]*EncounterRanking{
			encounter1.IDs[0]: encounterRanking1,
			encounter2.IDs[0]: encounterRanking2,
		},
	}
	encounters := &Encounters{Encounters: []*Encounter{encounter1, encounter2}}
	assert.Equal(t, bestRank, encounters.BestRankForEncounterRankings(encounterRankings))
}

func TestBestRankForEncounterRankings_NoClears(t *testing.T) {
	encounter1 := &Encounter{
		Name: "Encounter 1",
		IDs:  []int{1},
	}
	encounterRanking1 := &EncounterRanking{
		TotalKills: 0,
	}

	encounter2 := &Encounter{
		Name: "Encounter 2",
		IDs:  []int{2},
	}
	encounterRanking2 := &EncounterRanking{
		TotalKills: 0,
	}

	encounterRankings := &EncounterRankings{
		Encounters: map[int]*EncounterRanking{
			encounter1.IDs[0]: encounterRanking1,
			encounter2.IDs[0]: encounterRanking2,
		},
	}
	encounters := &Encounters{Encounters: []*Encounter{encounter1, encounter2}}
	assert.Equal(t, (*Rank)(nil), encounters.BestRankForEncounterRankings(encounterRankings))
}

func TestBestRankForEncounterRankings_BetterNotIncluded(t *testing.T) {
	encounter1 := &Encounter{
		Name: "Encounter 1",
		IDs:  []int{1},
	}
	bestRank := &Rank{Percent: 90.1384}
	encounterRanking1 := &EncounterRanking{
		TotalKills: 2,
		Ranks: []*Rank{
			bestRank,
			{Percent: 10.13},
			{Percent: 50.154},
		},
	}

	encounter2 := &Encounter{
		Name: "Encounter 2",
		IDs:  []int{2},
	}
	encounterRanking2 := &EncounterRanking{
		TotalKills: 3,
		Ranks: []*Rank{
			{Percent: 1.1},
			{Percent: 2.2},
		},
	}

	encounter3 := &Encounter{
		Name: "Encounter 3",
		IDs:  []int{3},
	}
	bestestRank := &Rank{Percent: 100.0}
	encounterRanking3 := &EncounterRanking{
		TotalKills: 3,
		Ranks: []*Rank{
			bestestRank,
			{Percent: 10.13},
			{Percent: 50.154},
		},
	}

	encounterRankings := &EncounterRankings{
		Encounters: map[int]*EncounterRanking{
			encounter1.IDs[0]: encounterRanking1,
			encounter2.IDs[0]: encounterRanking2,
			encounter3.IDs[0]: encounterRanking3,
		},
	}
	encounters := &Encounters{Encounters: []*Encounter{encounter1, encounter2}}
	assert.Equal(t, bestRank, encounters.BestRankForEncounterRankings(encounterRankings))
}

func TestTotalClearsFromEncounterRankings_SingleIDs(t *testing.T) {
	encounter1 := &Encounter{
		Name: "Encounter 1",
		IDs:  []int{1},
	}
	encounterRanking1 := &EncounterRanking{
		TotalKills: 2,
		Ranks: []*Rank{
			{Percent: 10.13},
			{Percent: 50.154},
		},
	}

	encounter2 := &Encounter{
		Name: "Encounter 2",
		IDs:  []int{2},
	}
	encounterRanking2 := &EncounterRanking{
		TotalKills: 3,
		Ranks: []*Rank{
			{Percent: 1.1},
			{Percent: 2.2},
		},
	}

	encounterRankings := &EncounterRankings{
		Encounters: map[int]*EncounterRanking{
			encounter1.IDs[0]: encounterRanking1,
			encounter2.IDs[0]: encounterRanking2,
		},
	}
	encounters := &Encounters{Encounters: []*Encounter{encounter1, encounter2}}
	assert.Equal(t, 2, encounters.TotalClearsFromEncounterRankings(encounterRankings))
}

func TestTotalClearsFromEncounterRankings_MultipleIDs(t *testing.T) {
	encounter1 := &Encounter{
		Name: "Encounter 1",
		IDs:  []int{1, 10, 100},
	}
	encounterRanking1 := &EncounterRanking{
		TotalKills: 2,
		Ranks: []*Rank{
			{Percent: 10.13},
			{Percent: 50.154},
		},
	}

	encounter2 := &Encounter{
		Name: "Encounter 2",
		IDs:  []int{2, 20, 200},
	}
	encounterRanking2 := &EncounterRanking{
		TotalKills: 3,
		Ranks: []*Rank{
			{Percent: 1.1},
			{Percent: 2.2},
		},
	}

	emptyRankings := &EncounterRanking{TotalKills: 0}

	encounterRankings := &EncounterRankings{
		Encounters: map[int]*EncounterRanking{
			encounter1.IDs[0]: encounterRanking1,
			encounter1.IDs[1]: emptyRankings,
			encounter1.IDs[2]: encounterRanking1,

			encounter2.IDs[0]: emptyRankings,
			encounter2.IDs[1]: emptyRankings,
			encounter2.IDs[2]: encounterRanking2,
		},
	}
	encounters := &Encounters{Encounters: []*Encounter{encounter1, encounter2}}
	assert.Equal(t, 2, encounters.TotalClearsFromEncounterRankings(encounterRankings))
}
