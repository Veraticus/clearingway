package clearingway

import (
	"fmt"
	"strings"

	"github.com/Veraticus/clearingway/internal/fflogs"
	"github.com/Veraticus/clearingway/internal/ffxiv"

	"github.com/bwmarrin/discordgo"
)

type ShouldApplyOpts struct {
	Character  *ffxiv.Character
	Encounters *Encounters
	Rankings   *fflogs.Rankings
}

type RoleType string

var (
	PfRole      RoleType = "PF"
	ReclearRole RoleType = "Reclear"
	ParseRole   RoleType = "Parse"
	ClearedRole RoleType = "Cleared"
)

type Roles struct {
	Roles []*Role
}
type Role struct {
	Type        RoleType
	Name        string
	Color       int
	ShouldApply func(*ShouldApplyOpts) (bool, string)
	DiscordRole *discordgo.Role
}

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

						for _, rank := range ranking.Ranks {
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
	}}
}

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
				"  `%v` with `%v` on <t:%v:F> (%v).\n",
				name,
				rank.Job.Abbreviation,
				rank.StartTime,
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
				percent := rank.Percent

				if percent < 1 {
					return true, fmt.Sprintf(
						"Parsed *0* (`%v`) with `%v` in `%v` on <t:%v:F>",
						rank.Percent,
						rank.Job.Abbreviation,
						encounter.Name,
						rank.StartTime,
					)
				}
				return false, "Worst parse was not 0."
			},
		},
	}}
}

func WorldRoles() *Roles {
	return &Roles{Roles: []*Role{
		{
			Name: "Aether", Color: 0x71368a,
			ShouldApply: func(opts *ShouldApplyOpts) (bool, string) {
				_, ok := ffxiv.AetherWorlds[opts.Character.World]
				if ok {
					return true, fmt.Sprintf("Charater `%v (%v)` is in Aether.", opts.Character.Name(), opts.Character.World)
				}
				return false, fmt.Sprintf("Character `%v (%v)` is not in Aether.", opts.Character.Name(), opts.Character.World)
			},
		},
		{
			Name: "Crystal", Color: 0x206694,
			ShouldApply: func(opts *ShouldApplyOpts) (bool, string) {
				_, ok := ffxiv.CrystalWorlds[opts.Character.World]
				if ok {
					return true, fmt.Sprintf("Charater `%v (%v)` is in Crystal.", opts.Character.Name(), opts.Character.World)
				}
				return false, fmt.Sprintf("Character `%v (%v)` is not in Crystal.", opts.Character.Name(), opts.Character.World)
			},
		},
		{
			Name: "Primal", Color: 0x992d22,
			ShouldApply: func(opts *ShouldApplyOpts) (bool, string) {
				_, ok := ffxiv.PrimalWorlds[opts.Character.World]
				if ok {
					return true, fmt.Sprintf("Charater `%v (%v)` is in Primal.", opts.Character.Name(), opts.Character.World)
				}
				return false, fmt.Sprintf("Character `%v (%v)` is not in Primal.", opts.Character.Name(), opts.Character.World)
			},
		},
	}}
}

func (r *Role) Ensure(guildId string, s *discordgo.Session, existingRoles []*discordgo.Role) error {
	var existingRole *discordgo.Role
	for _, er := range existingRoles {
		if er.Name == r.Name {
			existingRole = er
		}
	}
	if existingRole == nil {
		newRole, err := s.GuildRoleCreate(guildId)
		if err != nil {
			return fmt.Errorf("Could not create new role for %v: %w", r.Name, err)
		}
		existingRole = newRole
	}

	if existingRole.Color != r.Color || existingRole.Name != r.Name {
		newRole, err := s.GuildRoleEdit(
			guildId,
			existingRole.ID,
			r.Name,
			r.Color,
			false,
			0,
			false,
		)
		if err != nil {
			return fmt.Errorf("Could not ensure role %v: %w", r.Name, err)
		}
		existingRole = newRole
	}
	r.DiscordRole = existingRole

	return nil
}

func (r *Role) AddToCharacter(guildId, userId string, s *discordgo.Session, char *ffxiv.Character) error {
	return s.GuildMemberRoleAdd(guildId, userId, r.DiscordRole.ID)
}

func (r *Role) RemoveFromCharacter(guildId, userId string, s *discordgo.Session, char *ffxiv.Character) error {
	return s.GuildMemberRoleRemove(guildId, userId, r.DiscordRole.ID)
}

func (r *Role) PresentInRoles(existingRoleIds []string) bool {
	for _, existingRoleId := range existingRoleIds {
		if existingRoleId == r.DiscordRole.ID {
			return true
		}
	}

	return false
}

func (rs *Roles) Ensure(guildId string, s *discordgo.Session, existingRoles []*discordgo.Role) error {
	for _, r := range rs.Roles {
		err := r.Ensure(guildId, s, existingRoles)
		fmt.Printf("Ensuring role: %+v\n", r)
		if err != nil {
			return fmt.Errorf("Could not ensure role %v: %w", r, err)
		}
	}
	return nil
}

func (rs *Roles) Reorder(guildId string, s *discordgo.Session) error {
	discordRoles := []*discordgo.Role{}

	for _, role := range rs.Roles {
		discordRoles = append(discordRoles, role.DiscordRole)
	}

	_, err := s.GuildRoleReorder(guildId, discordRoles)
	return err
}
