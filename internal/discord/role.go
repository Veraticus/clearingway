package discord

import (
	"fmt"

	"github.com/Veraticus/clearingway/internal/fflogs"
	"github.com/Veraticus/clearingway/internal/ffxiv"

	"github.com/bwmarrin/discordgo"
)

var AetherServers = map[string]interface{}{
	"Adamantoise":  nil,
	"Cactuar":      nil,
	"Faerie":       nil,
	"Gilgamesh":    nil,
	"Jenova":       nil,
	"Midgardsormr": nil,
	"Sargatanas":   nil,
	"Siren":        nil,
}

var CrystalServers = map[string]interface{}{
	"Balmung":   nil,
	"Brynhildr": nil,
	"Coeurl":    nil,
	"Diabolos":  nil,
	"Goblin":    nil,
	"Malboro":   nil,
	"Mateus":    nil,
	"Zalera":    nil,
}

var PrimalServers = map[string]interface{}{
	"Behemoth":  nil,
	"Excalibur": nil,
	"Exodus":    nil,
	"Famfrit":   nil,
	"Hyperion":  nil,
	"Lamia":     nil,
	"Leviathan": nil,
	"Ultros":    nil,
}

type Roles struct {
	Roles []*Role
}

type Role struct {
	DiscordRole *discordgo.Role
	Name        string
	Color       int
	ShouldApply func(*ffxiv.Character, *fflogs.Encounters, *fflogs.EncounterRankings) bool
}

func AllParsingRoles() []*Role {
	return []*Role{
		{
			Name: "Gold", Color: 0xe1cc8a,
			ShouldApply: func(c *ffxiv.Character, es *fflogs.Encounters, ers *fflogs.EncounterRankings) bool {
				rank := es.BestRankForEncounterRankings(ers)
				if rank == nil {
					return false
				}
				percent := rank.Percent
				return percent == 100.0
			},
		},
		{
			Name: "Pink", Color: 0xd06fa4,
			ShouldApply: func(c *ffxiv.Character, es *fflogs.Encounters, ers *fflogs.EncounterRankings) bool {
				rank := es.BestRankForEncounterRankings(ers)
				if rank == nil {
					return false
				}
				percent := rank.Percent

				return (percent >= 99.0 && percent < 100.0)
			},
		},
		{
			Name: "Orange", Color: 0xef8633,
			ShouldApply: func(c *ffxiv.Character, es *fflogs.Encounters, ers *fflogs.EncounterRankings) bool {
				rank := es.BestRankForEncounterRankings(ers)
				if rank == nil {
					return false
				}
				percent := rank.Percent

				return (percent >= 95.0 && percent < 99.0)
			},
		},
		{
			Name: "Purple", Color: 0x9644e5,
			ShouldApply: func(c *ffxiv.Character, es *fflogs.Encounters, ers *fflogs.EncounterRankings) bool {
				rank := es.BestRankForEncounterRankings(ers)
				if rank == nil {
					return false
				}
				percent := rank.Percent

				return (percent >= 75.0 && percent < 95.0)
			},
		},
		{
			Name: "Blue", Color: 0x2a72f6,
			ShouldApply: func(c *ffxiv.Character, es *fflogs.Encounters, ers *fflogs.EncounterRankings) bool {
				rank := es.BestRankForEncounterRankings(ers)
				if rank == nil {
					return false
				}
				percent := rank.Percent

				return (percent >= 50.0 && percent < 75.0)
			},
		},
		{
			Name: "Green", Color: 0x78fa4c,
			ShouldApply: func(c *ffxiv.Character, es *fflogs.Encounters, ers *fflogs.EncounterRankings) bool {
				rank := es.BestRankForEncounterRankings(ers)
				if rank == nil {
					return false
				}
				percent := rank.Percent

				return (percent >= 25.0 && percent < 50.0)
			},
		},
		{
			Name: "Gray", Color: 0x636363,
			ShouldApply: func(c *ffxiv.Character, es *fflogs.Encounters, ers *fflogs.EncounterRankings) bool {
				rank := es.BestRankForEncounterRankings(ers)
				if rank == nil {
					return false
				}
				percent := rank.Percent

				return (percent > 0 && percent < 25.0)
			},
		},
		{
			Name: "NA's Comfiest", Color: 0x636363,
			ShouldApply: func(c *ffxiv.Character, es *fflogs.Encounters, ers *fflogs.EncounterRankings) bool {
				rank := es.BestRankForEncounterRankings(ers)
				if rank == nil {
					return false
				}
				percent := rank.Percent

				return percent == 0.0
			},
		},
	}
}

func AllUltimateRoles() []*Role {
	return []*Role{
		{
			Name: "The Legend", Color: 0x3498db,
			ShouldApply: func(c *ffxiv.Character, es *fflogs.Encounters, ers *fflogs.EncounterRankings) bool {
				totalClears := fflogs.UltimateEncounters.TotalClearsFromEncounterRankings(ers)
				return totalClears == 1
			},
		},
		{
			Name: "The Double Legend", Color: 0x3498db,
			ShouldApply: func(c *ffxiv.Character, es *fflogs.Encounters, ers *fflogs.EncounterRankings) bool {
				totalClears := fflogs.UltimateEncounters.TotalClearsFromEncounterRankings(ers)
				return totalClears == 2
			},
		},
		{
			Name: "The Triple Legend", Color: 0x3498db,
			ShouldApply: func(c *ffxiv.Character, es *fflogs.Encounters, ers *fflogs.EncounterRankings) bool {
				totalClears := fflogs.UltimateEncounters.TotalClearsFromEncounterRankings(ers)
				return totalClears == 3
			},
		},
		{
			Name: "The Quad Legend", Color: 0x3498db,
			ShouldApply: func(c *ffxiv.Character, es *fflogs.Encounters, ers *fflogs.EncounterRankings) bool {
				totalClears := fflogs.UltimateEncounters.TotalClearsFromEncounterRankings(ers)
				return totalClears == 4
			},
		},
	}
}

func AllServerRoles() []*Role {
	roles := []*Role{
		{
			Name: "Aether", Color: 0x71368a,
			ShouldApply: func(c *ffxiv.Character, es *fflogs.Encounters, ers *fflogs.EncounterRankings) bool {
				_, ok := AetherServers[c.World]
				return ok
			},
		},
		{
			Name: "Crystal", Color: 0x206694,
			ShouldApply: func(c *ffxiv.Character, es *fflogs.Encounters, ers *fflogs.EncounterRankings) bool {
				_, ok := CrystalServers[c.World]
				return ok
			},
		},
		{
			Name: "Primal", Color: 0x992d22,
			ShouldApply: func(c *ffxiv.Character, es *fflogs.Encounters, ers *fflogs.EncounterRankings) bool {
				_, ok := PrimalServers[c.World]
				return ok
			},
		},
	}

	return roles
}

func RolesForEncounters(es *fflogs.Encounters) []*Role {
	roles := []*Role{}

	for _, encounter := range es.Encounters {
		if encounter.CreateRoles != true {
			continue
		}

		encounterId := encounter.IDs[0]
		roles = append(roles, &Role{Name: encounter.Name + "-PF", Color: 0x11806a})
		roles = append(roles, &Role{Name: encounter.Name + "-Reclear", Color: 0x11806a})
		roles = append(roles, &Role{Name: encounter.Name + "-Parse", Color: 0x11806a})
		roles = append(roles, &Role{
			Name: encounter.Name + "-Cleared", Color: encounter.ClearedColor,
			ShouldApply: func(c *ffxiv.Character, es *fflogs.Encounters, ers *fflogs.EncounterRankings) bool {
				encounterRanking, ok := ers.Encounters[encounterId]
				if !ok {
					return false
				}

				return encounterRanking.Cleared()
			},
		})
	}

	return roles
}

func (rs *Roles) Ensure(guildId string, s *discordgo.Session, existingRoles []*discordgo.Role) error {
	for _, r := range rs.Roles {
		fmt.Printf("Ensuring %v...\n", r)
		err := r.Ensure(guildId, s, existingRoles)
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

func (rs *Roles) RoleNames(roleIds []string) []string {
	roleNames := []string{}

	for _, role := range rs.Roles {
		for _, roleId := range roleIds {
			if role.DiscordRole.ID == roleId {
				roleNames = append(roleNames, role.DiscordRole.Name)
			}
		}
	}

	return roleNames
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
