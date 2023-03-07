package clearingway

func RelevantParsingRoles() *Roles {
	return &Roles{Roles: []*Role{
		{
			Name: "Gold", Color: 0xe1cc8a, Uncolor: true,
			Description: "DPS parse is 100 in a relevant encounter.",
			ShouldApply: func(opts *ShouldApplyOpts) (bool, string) {
				encounter, rank := opts.Encounters.BestDPSRank(opts.Rankings)
				if encounter == nil || rank == nil {
					return false, "No encounter or rank found."
				}
				percent := rank.DPSPercent

				if percent == 100 {
					return true, rank.BestDPSParseString(encounter.Name)
				}
				return false, "Best parse was not 100."
			},
		},
		{
			Name: "Pink", Color: 0xd06fa4, Uncolor: true,
			Description: "DPS parse is between 99 and 100 in a relevant encounter.",
			ShouldApply: func(opts *ShouldApplyOpts) (bool, string) {
				encounter, rank := opts.Encounters.BestDPSRank(opts.Rankings)
				if encounter == nil || rank == nil {
					return false, "No encounter or rank found."
				}
				percent := rank.DPSPercent

				if percent >= 99.0 && percent < 100.0 {
					return true, rank.BestDPSParseString(encounter.Name)
				}
				return false, "Best parse was not between 99 and 100."
			},
		},
		{
			Name: "Orange", Color: 0xef8633, Uncolor: true,
			Description: "DPS parse is between 95 and 99 in a relevant encounter.",
			ShouldApply: func(opts *ShouldApplyOpts) (bool, string) {
				encounter, rank := opts.Encounters.BestDPSRank(opts.Rankings)
				if encounter == nil || rank == nil {
					return false, "No encounter or rank found."
				}
				percent := rank.DPSPercent

				if percent >= 95.0 && percent < 99.0 {
					return true, rank.BestDPSParseString(encounter.Name)
				}
				return false, "Best parse was not between 95 and 99."
			},
		},
		{
			Name: "Purple", Color: 0x9644e5, Uncolor: true,
			Description: "DPS parse is between 75 and 95 in a relevant encounter.",
			ShouldApply: func(opts *ShouldApplyOpts) (bool, string) {
				encounter, rank := opts.Encounters.BestDPSRank(opts.Rankings)
				if encounter == nil || rank == nil {
					return false, "No encounter or rank found."
				}
				percent := rank.DPSPercent

				if percent >= 75.0 && percent < 95.0 {
					return true, rank.BestDPSParseString(encounter.Name)
				}
				return false, "Best parse was not between 75 and 95."
			},
		},
		{
			Name: "Blue", Color: 0x2a72f6, Uncolor: true,
			Description: "DPS parse is between 50 and 75 in a relevant encounter.",
			ShouldApply: func(opts *ShouldApplyOpts) (bool, string) {
				encounter, rank := opts.Encounters.BestDPSRank(opts.Rankings)
				if encounter == nil || rank == nil {
					return false, "No encounter or rank found."
				}
				percent := rank.DPSPercent

				if percent >= 50.0 && percent < 75.0 {
					return true, rank.BestDPSParseString(encounter.Name)
				}
				return false, "Best parse was not between 50 and 75."
			},
		},
		{
			Name: "Green", Color: 0x78fa4c, Uncolor: true,
			Description: "DPS parse is between 25 and 50 in a relevant encounter.",
			ShouldApply: func(opts *ShouldApplyOpts) (bool, string) {
				encounter, rank := opts.Encounters.BestDPSRank(opts.Rankings)
				if encounter == nil || rank == nil {
					return false, "No encounter or rank found."
				}
				percent := rank.DPSPercent

				if percent >= 25.0 && percent < 50.0 {
					return true, rank.BestDPSParseString(encounter.Name)
				}
				return false, "Best parse was not between 25 and 50."
			},
		},
		{
			Name: "Gray", Color: 0x636363, Uncolor: true,
			Description: "DPS parse is between 0 and 25 in a relevant encounter.",
			ShouldApply: func(opts *ShouldApplyOpts) (bool, string) {
				encounter, rank := opts.Encounters.BestDPSRank(opts.Rankings)
				if encounter == nil || rank == nil {
					return false, "No encounter or rank found."
				}
				percent := rank.DPSPercent

				if percent > 0 && percent < 25.0 {
					return true, rank.BestDPSParseString(encounter.Name)
				}
				return false, "Best parse was not between 0 and 25."
			},
		},
	}}
}
