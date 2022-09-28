package clearingway

import (
	"fmt"
	"strings"

	"github.com/Veraticus/clearingway/internal/fflogs"
)

func ultimateRoleString(clearedEncounters *Encounters, rankings *fflogs.Rankings) string {
	clears := map[string]*fflogs.Ranking{}

	for _, clearedEncounter := range clearedEncounters.Encounters {
		for _, encounterId := range clearedEncounter.Ids {
			ranking, ok := rankings.Rankings[encounterId]
			if !ok {
				continue
			}
			if !ranking.Cleared() {
				continue
			}

			clears[clearedEncounter.Name] = ranking
		}
	}

	clearedString := strings.Builder{}
	clearedString.WriteString("Cleared the following Ultimate fights:\n")
	for name, ranking := range clears {
		rank := ranking.RanksByTime()[0]
		clearedString.WriteString(
			fmt.Sprintf(
				"     `%v` with `%v` on <t:%v:F> (%v).\n",
				name,
				rank.Job.Abbreviation,
				rank.UnixTime(),
				rank.Report.Url(),
			),
		)
	}

	return strings.TrimSuffix(clearedString.String(), "\n")
}

func UltimateRoles() *Roles {
	return &Roles{Roles: []*Role{
		{
			Name: "The Legend", Color: 0x3498db,
			ShouldApply: func(opts *ShouldApplyOpts) (bool, string) {
				clearedEncounters := opts.Encounters.Clears(opts.Rankings)
				if len(clearedEncounters.Encounters) == 1 {
					output := ultimateRoleString(clearedEncounters, opts.Rankings)
					return true, output
				}

				return false, "Did not clear only one ultimate."
			},
		},
		{
			Name: "The Double Legend", Color: 0x3498db,
			ShouldApply: func(opts *ShouldApplyOpts) (bool, string) {
				clearedEncounters := opts.Encounters.Clears(opts.Rankings)
				if len(clearedEncounters.Encounters) == 2 {
					output := ultimateRoleString(clearedEncounters, opts.Rankings)
					return true, output
				}

				return false, "Did not clear only two ultimates."
			},
		},
		{
			Name: "The Triple Legend", Color: 0x3498db,
			ShouldApply: func(opts *ShouldApplyOpts) (bool, string) {
				clearedEncounters := opts.Encounters.Clears(opts.Rankings)
				if len(clearedEncounters.Encounters) == 3 {
					output := ultimateRoleString(clearedEncounters, opts.Rankings)
					return true, output
				}

				return false, "Did not clear only three ultimates."
			},
		},
		{
			Name: "The Quad Legend", Color: 0x3498db,
			ShouldApply: func(opts *ShouldApplyOpts) (bool, string) {
				clearedEncounters := opts.Encounters.Clears(opts.Rankings)
				if len(clearedEncounters.Encounters) == 4 {
					output := ultimateRoleString(clearedEncounters, opts.Rankings)
					return true, output
				}

				return false, "Did not clear all four ultimates."
			},
		},
		{
			Name: "The Nice Legend", Color: 0xE48CA3,
			ShouldApply: func(opts *ShouldApplyOpts) (bool, string) {
				for _, encounter := range opts.Encounters.Encounters {
					for _, encounterId := range encounter.Ids {
						ranking, ok := opts.Rankings.Rankings[encounterId]
						if !ok {
							continue
						}
						if !ranking.Cleared() {
							continue
						}

						for _, rank := range ranking.Ranks {
							if rank.DPSPercent >= 69.0 && rank.DPSPercent <= 69.9 {
								return true,
									fmt.Sprintf(
										"Parsed **69** (`%v`) with `%v` in `%v` on <t:%v:F> (%v).",
										rank.DPSPercentString(),
										rank.Job.Abbreviation,
										encounter.Name,
										rank.UnixTime(),
										rank.Report.Url(),
									)
							}
						}
					}
				}

				return false, "No ultimate encounter had a parse at 69."
			},
		},
		{
			Name: "The Comfy Legend", Color: 0x636363,
			ShouldApply: func(opts *ShouldApplyOpts) (bool, string) {
				encounter, rank := opts.Encounters.WorstDPSRank(opts.Rankings)
				if encounter == nil || rank == nil {
					return false, "No encounter or rank found."
				}
				percent := rank.DPSPercent

				if rank.DPSParseFound && percent < 1 {
					return true, fmt.Sprintf(
						"Parsed **0** (`%v`) with `%v` in `%v` on <t:%v:F> (%v).",
						rank.DPSPercentString(),
						rank.Job.Abbreviation,
						encounter.Name,
						rank.UnixTime(),
						rank.Report.Url(),
					)
				}
				return false, "No ultimate encounter had a parse at 0."
			},
		},
		{
			Name: "The Chadding Legend", Color: 0x39FF14,
			ShouldApply: func(opts *ShouldApplyOpts) (bool, string) {
				for _, encounter := range opts.Encounters.Encounters {
					for _, encounterId := range encounter.Ids {
						ranking, ok := opts.Rankings.Rankings[encounterId]
						if !ok {
							continue
						}
						if !ranking.Cleared() {
							continue
						}

						for _, rank := range ranking.Ranks {
							if rank.HPSParseFound && rank.HPSPercent < 1 && rank.Job.IsHealer() {
								return true,
									fmt.Sprintf(
										"HPS parsed was **0** (`%v`) as a healer (`%v`) in `%v` on <t:%v:F> (%v).",
										rank.HPSPercentString(),
										rank.Job.Abbreviation,
										encounter.Name,
										rank.UnixTime(),
										rank.Report.Url(),
									)
							}
						}
					}
				}

				return false, "No ultimate encounter had a HPS parse at 0."
			},
		},
		{
			Name: "The Bloodbathing Legend", Color: 0x8a0303,
			ShouldApply: func(opts *ShouldApplyOpts) (bool, string) {
				for _, encounter := range opts.Encounters.Encounters {
					for _, encounterId := range encounter.Ids {
						ranking, ok := opts.Rankings.Rankings[encounterId]
						if !ok {
							continue
						}
						if !ranking.Cleared() {
							continue
						}

						for _, rank := range ranking.Ranks {
							if rank.HPSParseFound && rank.HPSPercent == 100 && !rank.Job.IsHealer() {
								return true,
									fmt.Sprintf(
										"HPS parsed was **100** (`%v`) as a non-healer (`%v`) in `%v` on <t:%v:F> (%v).",
										rank.HPSPercentString(),
										rank.Job.Abbreviation,
										encounter.Name,
										rank.UnixTime(),
										rank.Report.Url(),
									)
							}
						}
					}
				}

				return false, "No encounter had a non-healer HPS parse at 100."
			},
		},
		{
			Name: "The Overhealing Legend", Color: 0xFFFFFF,
			ShouldApply: func(opts *ShouldApplyOpts) (bool, string) {
				for _, encounter := range opts.Encounters.Encounters {
					for _, encounterId := range encounter.Ids {
						ranking, ok := opts.Rankings.Rankings[encounterId]
						if !ok {
							continue
						}
						if !ranking.Cleared() {
							continue
						}

						for _, rank := range ranking.Ranks {
							if rank.HPSParseFound && rank.HPSPercent == 100 && rank.Job.IsHealer() {
								return true,
									fmt.Sprintf(
										"HPS parsed was **100** (`%v`) as a healer (`%v`) in `%v` on <t:%v:F> (%v).",
										rank.HPSPercentString(),
										rank.Job.Abbreviation,
										encounter.Name,
										rank.UnixTime(),
										rank.Report.Url(),
									)
							}
						}
					}
				}

				return false, "No encounter had a healer HPS parse at 100."
			},
		},
	}}
}
