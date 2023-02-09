package clearingway

import (
	"fmt"

	"github.com/Veraticus/clearingway/internal/fflogs"
	"github.com/Veraticus/clearingway/internal/ffxiv"

	"github.com/bwmarrin/discordgo"
)

type ShouldApplyOpts struct {
	Character     *ffxiv.Character
	Encounters    *Encounters
	Rankings      *fflogs.Rankings
	Fights        *fflogs.Fights
	ExistingRoles *Roles
}

type RoleType string

var (
	PfRole      RoleType = "PF"
	ReclearRole RoleType = "Reclear"
	ParseRole   RoleType = "Parse"
	ClearedRole RoleType = "Cleared"
	ProgRole    RoleType = "Prog"
)

type Roles struct {
	Roles []*Role

	// The first list of roles is roles to add; the second is roles to remove.
	ShouldApply func(*ShouldApplyOpts) (bool, string, []*Role, []*Role)
}

type Role struct {
	Type        RoleType
	Name        string
	Description string
	Color       int
	Uncomfy     bool
	Skip        bool
	Mention     bool
	Hoist       bool
	ShouldApply func(*ShouldApplyOpts) (bool, string)
	DiscordRole *discordgo.Role
}

func (r *Role) Ensure(guildId string, s *discordgo.Session, existingRoles []*discordgo.Role) error {
	if r.Skip {
		return nil
	}

	var existingRole *discordgo.Role
	for _, er := range existingRoles {
		if er.Name == r.Name {
			existingRole = er
		}
	}

	if existingRole == nil {
		newRole, err := s.GuildRoleCreate(guildId)
		if err != nil {
			return fmt.Errorf("Could not create new role for %v: %w.\n", r.Name, err)
		}
		existingRole = newRole
	}

	if existingRole.Color != r.Color || existingRole.Name != r.Name {
		newRole, err := s.GuildRoleEdit(
			guildId,
			existingRole.ID,
			r.Name,
			r.Color,
			r.Hoist,
			0,
			r.Mention,
		)
		if err != nil {
			return fmt.Errorf("Could not ensure role %v: %w.\n", r.Name, err)
		}
		existingRole = newRole
	}
	r.DiscordRole = existingRole

	return nil
}

func (r *Role) AddToCharacter(guildId, userId string, s *discordgo.Session) error {
	return s.GuildMemberRoleAdd(guildId, userId, r.DiscordRole.ID)
}

func (r *Role) RemoveFromCharacter(guildId, userId string, s *discordgo.Session) error {
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

func (r *Role) Phase(i int) string {
	if r.Type == ClearedRole {
		return "**Cleared**"
	}

	return fmt.Sprintf("Phase **%d**", i)
}

func (rs *Roles) Reorder(guildId string, s *discordgo.Session) error {
	discordRoles := []*discordgo.Role{}

	for _, role := range rs.Roles {
		discordRoles = append(discordRoles, role.DiscordRole)
	}

	_, err := s.GuildRoleReorder(guildId, discordRoles)
	return err
}

func (rs *Roles) FindByName(name string) *Role {
	for _, r := range rs.Roles {
		if r.Name == name {
			return r
		}
	}

	return nil
}

func (rs *Roles) IndexOfRole(role *Role) (bool, int) {
	for i, r := range rs.Roles {
		if r.Name == role.Name {
			return true, i
		}
	}

	return false, 0
}

func (rs *Roles) InDiscordRoles(ids []string) []*Role {
	roles := []*Role{}
	for _, r := range rs.Roles {
		for _, id := range ids {
			if r.DiscordRole.ID == id {
				roles = append(roles, r)
			}
		}
	}

	return roles
}
