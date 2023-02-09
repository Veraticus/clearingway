package fflogs

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Veraticus/clearingway/internal/ffxiv"
)

type report struct {
	Fights     []*fight    `json:"fights"`
	MasterData *masterData `json:"masterData"`
}

type fight struct {
	ID                       int   `json:"id"`
	Kill                     bool  `json:"kill"`
	Difficulty               int   `json:"difficulty"`
	EncounterID              int   `json:"encounterID"`
	LastPhaseAsAbsoluteIndex int   `json:"lastPhaseAsAbsoluteIndex"`
	FriendlyPlayers          []int `json:"friendlyPlayers"`
}

type masterData struct {
	Actors []*actor `json:"actors"`
}

type actor struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Server string `json:"server"`
}

type Fights struct {
	Fights []*Fight
}

type Fight struct {
	LastPhaseIndex int
	Kill           bool
	EncounterID    int
	ReportID       string
	ID             int
}

func (f *Fflogs) GetProgForReport(r string, rankingsToGet []*RankingToGet, char *ffxiv.Character) (*Fights, error) {
	query := strings.Builder{}
	query.WriteString(
		fmt.Sprintf(
			"query{reportData{report(code: \"%s\") {fights {kill difficulty id encounterID lastPhaseAsAbsoluteIndex friendlyPlayers} masterData(translate: false) {actors(type: \"Player\") {id name server}}}}}",
			r,
		),
	)

	raw, err := f.graphqlClient.ExecRaw(context.Background(), query.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("Error executing query: %w", err)
	}

	var response map[string]*json.RawMessage
	err = json.Unmarshal(raw, &response)
	if err != nil {
		return nil, fmt.Errorf("Could not unmarshal JSON: %w", err)
	}

	var reportData map[string]*json.RawMessage
	err = json.Unmarshal(*response["reportData"], &reportData)
	if err != nil {
		return nil, fmt.Errorf("Could not unmarshal JSON: %w", err)
	}

	var report *report
	err = json.Unmarshal(*reportData["report"], &report)
	if err != nil {
		return nil, fmt.Errorf("Could not unmarshal JSON: %w", err)
	}
	if report.Fights == nil {
		return nil, fmt.Errorf("Fight data not found correctly for %s!", r)
	}
	if report.MasterData == nil {
		return nil, fmt.Errorf("Master data not found correctly for %s!", r)
	}

	var characterActorId int
	characterFoundInMasterData := false
	for _, a := range report.MasterData.Actors {
		if a.Name == char.Name() && a.Server == char.World {
			characterActorId = a.Id
			characterFoundInMasterData = true
		}
	}
	if !characterFoundInMasterData {
		return nil, fmt.Errorf("Could not find character %s (%s) in report %s.", char.Name(), char.World, r)
	}

	fights := &Fights{Fights: []*Fight{}}

	for _, f := range report.Fights {
		for _, rankingToGet := range rankingsToGet {
			for _, encounterID := range rankingToGet.IDs {
				if rankingToGet.Difficulty != f.Difficulty {
					continue
				}

				if f.EncounterID != encounterID {
					continue
				}

				characterIsInFight := false
				for _, friendlyId := range f.FriendlyPlayers {
					if friendlyId == characterActorId {
						characterIsInFight = true
					}
				}
				if !characterIsInFight {
					continue
				}

				fights.Add(&Fight{
					Kill:           f.Kill,
					LastPhaseIndex: f.LastPhaseAsAbsoluteIndex,
					EncounterID:    f.EncounterID,
					ID:             f.ID,
					ReportID:       r,
				})
			}
		}
	}

	return fights, nil
}

func (f *Fights) Add(fight *Fight) {
	for _, existingFight := range f.Fights {
		if existingFight.ID == fight.ID && existingFight.ReportID == fight.ReportID {
			return
		}
	}

	f.Fights = append(f.Fights, fight)
}

func (f *Fights) FurthestFight() *Fight {
	var fight *Fight

	for _, f := range f.Fights {
		if f.Kill {
			return f
		}
		if fight == nil || f.LastPhaseIndex > fight.LastPhaseIndex {
			fight = f
		}
	}

	return fight
}

func (f *Fight) ReportURL() string {
	return fmt.Sprintf("https://www.fflogs.com/reports/%s#fight=%d", f.ReportID, f.ID)
}
