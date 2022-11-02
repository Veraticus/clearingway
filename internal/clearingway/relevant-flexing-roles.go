package clearingway

import (
	"fmt"
)

func RelevantFlexingRoles() *Roles {
	return &Roles{Roles: []*Role{
		{
			Name: "NA's Comfiest", Color: 0x636363, Uncomfy: true,
			Description: "DPS parse rounds to zero in a relevant encounter.",
			ShouldApply: func(opts *ShouldApplyOpts) (bool, string) {
				encounter, rank := opts.Encounters.WorstDPSRank(opts.Rankings)
				if encounter == nil || rank == nil {
					return false, "No encounter or rank found."
				}
				percent := rank.DPSPercent

				if rank.DPSParseFound && percent < 1 {
					return true, fmt.Sprintf(
						"Parsed **0** (%v) with `%v` in `%v` on <t:%v:F> (%v).\nUse `/uncomfy` if you don't want this role.",
						rank.DPSPercentString(),
						rank.Job.Abbreviation,
						encounter.Name,
						rank.UnixTime(),
						rank.Report.Url(),
					)
				}
				return false, "Worst parse was not 0."
			},
		},
		{
			Name: "Nice", Color: 0xE48CA3,
			Description: "DPS parse rounds to 69 (nice) in a relevant encounter.",
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
							if rank.DPSPercent >= 69.0 && rank.DPSPercent < 70.0 {
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

				return false, "No encounter had a parse at 69."
			},
		},
		{
			Name: "Chad", Color: 0x39FF14, Uncomfy: true,
			Description: "HPS parse as a healer rounds to 0 in a relevant encounter.",
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
										"HPS parsed was **0** (`%v`) as a healer (`%v`) in `%v` on <t:%v:F> (%v).\nUse `/uncomfy` if you don't want this role.",
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

				return false, "No encounter had a healer HPS parse at 0."
			},
		},
		{
			Name: "Bloodbather", Color: 0x8a0303,
			Description: "HPS parse as a non-healer is 100 in a relevant encounter.",
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
			Name: "Overhealer", Color: 0xFFFFFF,
			Description: "HPS parse as a healer is 100 in a relevant encounter.",
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
