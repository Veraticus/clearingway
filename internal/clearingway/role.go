package clearingway

import (
	"fmt"

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
	ShouldApply func(*ShouldApplyOpts) bool
	DiscordRole *discordgo.Role
}

func ParsingRoles() *Roles {
	return &Roles{Roles: []*Role{
		{
			Name: "Gold", Color: 0xe1cc8a,
			ShouldApply: func(opts *ShouldApplyOpts) bool {
				rank := opts.Encounters.BestRank(opts.Rankings)
				if rank == nil {
					return false
				}
				percent := rank.Percent
				return percent == 100.0
			},
		},
		{
			Name: "Pink", Color: 0xd06fa4,
			ShouldApply: func(opts *ShouldApplyOpts) bool {
				rank := opts.Encounters.BestRank(opts.Rankings)
				if rank == nil {
					return false
				}
				percent := rank.Percent

				return (percent >= 99.0 && percent < 100.0)
			},
		},
		{
			Name: "Orange", Color: 0xef8633,
			ShouldApply: func(opts *ShouldApplyOpts) bool {
				rank := opts.Encounters.BestRank(opts.Rankings)
				if rank == nil {
					return false
				}
				percent := rank.Percent

				return (percent >= 95.0 && percent < 99.0)
			},
		},
		{
			Name: "Purple", Color: 0x9644e5,
			ShouldApply: func(opts *ShouldApplyOpts) bool {
				rank := opts.Encounters.BestRank(opts.Rankings)
				if rank == nil {
					return false
				}
				percent := rank.Percent

				return (percent >= 75.0 && percent < 95.0)
			},
		},
		{
			Name: "Blue", Color: 0x2a72f6,
			ShouldApply: func(opts *ShouldApplyOpts) bool {
				rank := opts.Encounters.BestRank(opts.Rankings)
				if rank == nil {
					return false
				}
				percent := rank.Percent

				return (percent >= 50.0 && percent < 75.0)
			},
		},
		{
			Name: "Green", Color: 0x78fa4c,
			ShouldApply: func(opts *ShouldApplyOpts) bool {
				rank := opts.Encounters.BestRank(opts.Rankings)
				if rank == nil {
					return false
				}
				percent := rank.Percent

				return (percent >= 25.0 && percent < 50.0)
			},
		},
		{
			Name: "Gray", Color: 0x636363,
			ShouldApply: func(opts *ShouldApplyOpts) bool {
				rank := opts.Encounters.BestRank(opts.Rankings)
				if rank == nil {
					return false
				}
				percent := rank.Percent

				return (percent > 0 && percent < 25.0)
			},
		},
		{
			Name: "NA's Comfiest", Color: 0x636363,
			ShouldApply: func(opts *ShouldApplyOpts) bool {
				rank := opts.Encounters.WorstRank(opts.Rankings)
				fmt.Printf("Worst rank is: %+v\n", rank)
				if rank == nil {
					return false
				}
				percent := rank.Percent

				return percent < 1
			},
		},
		{
			Name: "Nice", Color: 0xE48CA3,
			ShouldApply: func(opts *ShouldApplyOpts) bool {
				for _, encounter := range opts.Encounters.Encounters {
					for _, encounterId := range encounter.Ids {
						encounterRanking, ok := opts.Rankings.Rankings[encounterId]
						if !ok {
							continue
						}
						if !encounterRanking.Cleared() {
							continue
						}

						for _, rank := range encounterRanking.Ranks {
							if rank.Percent >= 69.0 && rank.Percent <= 69.9 {
								return true
							}
						}
					}
				}

				return false
			},
		},
	}}
}

func UltimateRoles() *Roles {
	return &Roles{Roles: []*Role{
		{
			Name: "The Legend", Color: 0x3498db,
			ShouldApply: func(opts *ShouldApplyOpts) bool {
				totalClears := opts.Encounters.TotalClears(opts.Rankings)
				return totalClears == 1
			},
		},
		{
			Name: "The Double Legend", Color: 0x3498db,
			ShouldApply: func(opts *ShouldApplyOpts) bool {
				totalClears := opts.Encounters.TotalClears(opts.Rankings)
				return totalClears == 2
			},
		},
		{
			Name: "The Triple Legend", Color: 0x3498db,
			ShouldApply: func(opts *ShouldApplyOpts) bool {
				totalClears := opts.Encounters.TotalClears(opts.Rankings)
				return totalClears == 3
			},
		},
		{
			Name: "The Quad Legend", Color: 0x3498db,
			ShouldApply: func(opts *ShouldApplyOpts) bool {
				totalClears := opts.Encounters.TotalClears(opts.Rankings)
				return totalClears == 4
			},
		},
		{
			Name: "The Nice Legend", Color: 0xE48CA3,
			ShouldApply: func(opts *ShouldApplyOpts) bool {
				for _, encounter := range opts.Encounters.Encounters {
					for _, encounterId := range encounter.Ids {
						encounterRanking, ok := opts.Rankings.Rankings[encounterId]
						if !ok {
							continue
						}
						if !encounterRanking.Cleared() {
							continue
						}

						for _, rank := range encounterRanking.Ranks {
							if rank.Percent >= 69.0 && rank.Percent <= 69.9 {
								return true
							}
						}
					}
				}

				return false
			},
		},
		{
			Name: "The Comfy Legend", Color: 0x636363,
			ShouldApply: func(opts *ShouldApplyOpts) bool {
				rank := opts.Encounters.WorstRank(opts.Rankings)
				if rank == nil {
					return false
				}
				percent := rank.Percent

				return percent <= 0.9
			},
		},
	}}
}

func WorldRoles() *Roles {
	return &Roles{Roles: []*Role{
		{
			Name: "Aether", Color: 0x71368a,
			ShouldApply: func(opts *ShouldApplyOpts) bool {
				_, ok := ffxiv.AetherWorlds[opts.Character.World]
				return ok
			},
		},
		{
			Name: "Crystal", Color: 0x206694,
			ShouldApply: func(opts *ShouldApplyOpts) bool {
				_, ok := ffxiv.CrystalWorlds[opts.Character.World]
				return ok
			},
		},
		{
			Name: "Primal", Color: 0x992d22,
			ShouldApply: func(opts *ShouldApplyOpts) bool {
				_, ok := ffxiv.PrimalWorlds[opts.Character.World]
				return ok
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
