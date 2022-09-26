package clearingway

import (
	"fmt"
)

func ParsingRoles() *Roles {
	return &Roles{Roles: []*Role{
		{
			Name: "Gold", Color: 0xe1cc8a,
			ShouldApply: func(opts *ShouldApplyOpts) (bool, string) {
				encounter, rank := opts.Encounters.BestDPSRank(opts.Rankings)
				if encounter == nil || rank == nil {
					return false, "No encounter or rank found."
				}
				percent := rank.Percent

				if percent == 100 {
					return true, fmt.Sprintf(
						"Best parse was %v with %v in %v on <t:%v:F> (%v).",
						rank.Percent,
						rank.Job.Abbreviation,
						encounter.Name,
						rank.StartTime,
						rank.Report.Url(),
					)
				}
				return false, "Best parse was not 100."
			},
		},
		{
			Name: "Pink", Color: 0xd06fa4,
			ShouldApply: func(opts *ShouldApplyOpts) (bool, string) {
				encounter, rank := opts.Encounters.BestDPSRank(opts.Rankings)
				if encounter == nil || rank == nil {
					return false, "No encounter or rank found."
				}
				percent := rank.Percent

				if percent >= 99.0 && percent < 100.0 {
					return true, rank.BestParseString(encounter.Name)
				}
				return false, "Best parse was not between 99 and 100."
			},
		},
		{
			Name: "Orange", Color: 0xef8633,
			ShouldApply: func(opts *ShouldApplyOpts) (bool, string) {
				encounter, rank := opts.Encounters.BestDPSRank(opts.Rankings)
				if encounter == nil || rank == nil {
					return false, "No encounter or rank found."
				}
				percent := rank.Percent

				if percent >= 95.0 && percent < 99.0 {
					return true, rank.BestParseString(encounter.Name)
				}
				return false, "Best parse was not between 95 and 99."
			},
		},
		{
			Name: "Purple", Color: 0x9644e5,
			ShouldApply: func(opts *ShouldApplyOpts) (bool, string) {
				encounter, rank := opts.Encounters.BestDPSRank(opts.Rankings)
				if encounter == nil || rank == nil {
					return false, "No encounter or rank found."
				}
				percent := rank.Percent

				if percent >= 75.0 && percent < 95.0 {
					return true, rank.BestParseString(encounter.Name)
				}
				return false, "Best parse was not between 75 and 95."
			},
		},
		{
			Name: "Blue", Color: 0x2a72f6,
			ShouldApply: func(opts *ShouldApplyOpts) (bool, string) {
				encounter, rank := opts.Encounters.BestDPSRank(opts.Rankings)
				if encounter == nil || rank == nil {
					return false, "No encounter or rank found."
				}
				percent := rank.Percent

				if percent >= 50.0 && percent < 75.0 {
					return true, rank.BestParseString(encounter.Name)
				}
				return false, "Best parse was not between 50 and 75."
			},
		},
		{
			Name: "Green", Color: 0x78fa4c,
			ShouldApply: func(opts *ShouldApplyOpts) (bool, string) {
				encounter, rank := opts.Encounters.BestDPSRank(opts.Rankings)
				if encounter == nil || rank == nil {
					return false, "No encounter or rank found."
				}
				percent := rank.Percent

				if percent >= 25.0 && percent < 50.0 {
					return true, rank.BestParseString(encounter.Name)
				}
				return false, "Best parse was not between 25 and 50."
			},
		},
		{
			Name: "Gray", Color: 0x636363,
			ShouldApply: func(opts *ShouldApplyOpts) (bool, string) {
				encounter, rank := opts.Encounters.BestDPSRank(opts.Rankings)
				if encounter == nil || rank == nil {
					return false, "No encounter or rank found."
				}
				percent := rank.Percent

				if percent > 0 && percent < 25.0 {
					return true, rank.BestParseString(encounter.Name)
				}
				return false, "Best parse was not between 0 and 25."
			},
		},
		{
			Name: "NA's Comfiest", Color: 0x636363,
			ShouldApply: func(opts *ShouldApplyOpts) (bool, string) {
				encounter, rank := opts.Encounters.WorstDPSRank(opts.Rankings)
				if encounter == nil || rank == nil {
					return false, "No encounter or rank found."
				}
				percent := rank.Percent

				if percent < 1 {
					return true, fmt.Sprintf(
						"Parsed 0 (%v) with %v in %v on <t:%v:F>",
						rank.Percent,
						rank.Job.Abbreviation,
						encounter.Name,
						rank.StartTime,
					)
				}
				return false, "Worst parse was not 0."
			},
		},
		{
			Name: "Nice", Color: 0xE48CA3,
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

						for _, rank := range ranking.DPSRanks() {
							if rank.Percent >= 69.0 && rank.Percent <= 69.9 {
								return true,
									fmt.Sprintf(
										"Parsed *69* (`%v`) with `%v` in `%v` on <t:%v:F>",
										rank.Percent,
										rank.Job.Abbreviation,
										encounter.Name,
										rank.StartTime,
									)
							}
						}
					}
				}

				return false, "No encounter had a parse at 69."
			},
		},
		{
			Name: "Chad", Color: 0x39FF14,
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

						for _, rank := range ranking.HPSRanks() {
							if rank.Percent < 1 && rank.Job.IsHealer() {
								return true,
									fmt.Sprintf(
										"HPS parsed was *0* (`%v`) as a healer (`%v`) in `%v` on <t:%v:F>",
										rank.Percent,
										rank.Job.Abbreviation,
										encounter.Name,
										rank.StartTime,
									)
							}
						}
					}
				}

				return false, "No encounter had a HPS parse at 0."
			},
		},
	}}
}
