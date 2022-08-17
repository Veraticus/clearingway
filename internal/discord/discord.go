package discord

import (
	"fmt"

	"github.com/Veraticus/clearingway/internal/fflogs"
	"github.com/Veraticus/clearingway/internal/ffxiv"

	"github.com/bwmarrin/discordgo"
)

type Discord struct {
	Token              string
	ChannelId          string
	Fflogs             *fflogs.Fflogs
	RelevantEncounters *fflogs.Encounters

	Roles map[string]*discordgo.Role

	Session *discordgo.Session
}

func (d *Discord) Start() error {
	s, err := discordgo.New("Bot " + d.Token)
	if err != nil {
		return fmt.Errorf("Could not start Discord: %f", err)
	}

	s.AddHandler(d.ready)
	s.AddHandler(d.messageCreate)
	s.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages

	err = s.Open()
	if err != nil {
		return fmt.Errorf("Could not open Discord session: %f", err)
	}

	d.Session = s
	return nil
}

func (d *Discord) ready(s *discordgo.Session, event *discordgo.Ready) {
	fmt.Printf("Initializing Discord...\n")

	roles := []*discordgo.Role{}
	for _, guild := range event.Guilds {
		fmt.Printf("Ensuring parse roles...\n")
		role, err := d.ensureRole(guild.ID, "Gold", 0xe1cc8a)
		if err != nil {
			fmt.Printf("Could not create role Gold: %v", err)
			return
		}
		roles = append(roles, role)
		role, err = d.ensureRole(guild.ID, "Pink", 0xd06fa4)
		if err != nil {
			fmt.Printf("Could not create role Pink: %v", err)
			return
		}
		roles = append(roles, role)
		role, err = d.ensureRole(guild.ID, "Orange", 0xef8633)
		if err != nil {
			fmt.Printf("Could not create role Orange: %v", err)
			return
		}
		roles = append(roles, role)
		role, err = d.ensureRole(guild.ID, "Purple", 0x9644e5)
		if err != nil {
			fmt.Printf("Could not create role Purple: %v", err)
			return
		}
		roles = append(roles, role)
		role, err = d.ensureRole(guild.ID, "Blue", 0x2a72f6)
		if err != nil {
			fmt.Printf("Could not create role Blue: %v", err)
			return
		}
		roles = append(roles, role)
		role, err = d.ensureRole(guild.ID, "Green", 0x78fa4c)
		if err != nil {
			fmt.Printf("Could not create role Green: %v", err)
			return
		}
		roles = append(roles, role)
		role, err = d.ensureRole(guild.ID, "Gray", 0x636363)
		if err != nil {
			fmt.Printf("Could not create role Gray: %v", err)
			return
		}
		roles = append(roles, role)
		role, err = d.ensureRole(guild.ID, "NA's Comfiest", 0x636363)
		if err != nil {
			fmt.Printf("Could not create role NA's Comfiest: %v", err)
			return
		}
		roles = append(roles, role)

		fmt.Printf("Ensuring ultimate roles...\n")
		role, err = d.ensureRole(guild.ID, "The Legend", 0x3498db)
		if err != nil {
			fmt.Printf("Could not create role The Legend: %v", err)
			return
		}
		roles = append(roles, role)
		role, err = d.ensureRole(guild.ID, "The Double Legend", 0x3498db)
		if err != nil {
			fmt.Printf("Could not create role The Double Legend: %v", err)
			return
		}
		roles = append(roles, role)
		role, err = d.ensureRole(guild.ID, "The Triple Legend", 0x3498db)
		if err != nil {
			fmt.Printf("Could not create role The Triple Legend: %v", err)
			return
		}
		roles = append(roles, role)
		role, err = d.ensureRole(guild.ID, "The Tetra Legend", 0x3498db)
		if err != nil {
			fmt.Printf("Could not create role The Tetra Legend: %v", err)
			return
		}
		roles = append(roles, role)

		fmt.Printf("Ensuring datacenter roles...\n")
		role, err = d.ensureRole(guild.ID, "Aether", 0x71368a)
		if err != nil {
			fmt.Printf("Could not create role Aether: %v", err)
			return
		}
		roles = append(roles, role)
		for _, aetherServer := range ffxiv.AetherServers {
			role, err = d.ensureRole(guild.ID, aetherServer, 0x71368a)
			if err != nil {
				fmt.Printf("Could not create role %v: %v", aetherServer, err)
				return
			}
			roles = append(roles, role)
		}
		d.ensureRole(guild.ID, "Crystal", 0x206694)
		if err != nil {
			fmt.Printf("Could not create role Crystal: %v", err)
			return
		}
		roles = append(roles, role)
		for _, crystalServer := range ffxiv.CrystalServers {
			role, err = d.ensureRole(guild.ID, crystalServer, 0x206694)
			if err != nil {
				fmt.Printf("Could not create role %v: %v", crystalServer, err)
				return
			}
			roles = append(roles, role)

		}
		role, err = d.ensureRole(guild.ID, "Primal", 0x992d22)
		if err != nil {
			fmt.Printf("Could not create role Primal: %v", err)
			return
		}
		roles = append(roles, role)
		for _, primalServer := range ffxiv.PrimalServers {
			role, err = d.ensureRole(guild.ID, primalServer, 0x992d22)
			if err != nil {
				fmt.Printf("Could not create role %v: %v", primalServer, err)
				return
			}
		}

		fmt.Printf("Ensuring relevant encounter roles...\n")
		for _, encounter := range d.RelevantEncounters.Encounters {
			role, err = d.ensureRole(guild.ID, encounter.Name+"-PF", 0x11806a)
			if err != nil {
				fmt.Printf("Could not create role %v-PF: %v", encounter.Name, err)
				return
			}
			roles = append(roles, role)

			role, err = d.ensureRole(guild.ID, encounter.Name+"-Cleared", 0x11806a)
			if err != nil {
				fmt.Printf("Could not create role %v-Cleared: %v", encounter.Name, err)
				return
			}
			roles = append(roles, role)

		}

		fmt.Printf("Reorder roles...\n")
		s.GuildRoleReorder(guild.ID, roles)
	}
	fmt.Printf("Discord ready!\n")
}

func (d *Discord) ensureRole(guildId, name string, color int) (*discordgo.Role, error) {
	existingRoles, err := d.Session.GuildRoles(guildId)
	if err != nil {
		return nil, fmt.Errorf("Could not obtain guild roles: %v", err)
	}

	id := ""
	for _, existingRole := range existingRoles {
		if existingRole.Name == name {
			id = existingRole.ID
		}
	}
	if id == "" {
		newRole, err := d.Session.GuildRoleCreate(guildId)
		if err != nil {
			return nil, fmt.Errorf("Could not create new role for %v: %w", name, err)
		}
		id = newRole.ID
	}
	role, err := d.Session.GuildRoleEdit(
		guildId,
		id,
		name,
		color,
		false,
		0,
		false,
	)
	if err != nil {
		return nil, fmt.Errorf("Could not ensure role %v: %w", name, err)
	}
	return role, nil
}

func (d *Discord) messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore messages from the bot
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Ignore messages not on the correct channel
	if m.ChannelID != d.ChannelId {
		return
	}

	// Check if the message is "!clears"
	if m.Content == "!clears" {
		s.ChannelMessageSendReply(d.ChannelId, "Received message!", (*m).Reference())
		// s.GuildMemberRoleAdd()
		// s.GuildMemberRoleRemove()
	}
}

func (d *Discord) resetUser(guildId, userId string) {
	removeRoles := []string{"NA'S Comfiest", "Gray", "Green", "Blue", "Purple", "Orange", "Pink", "Gold", "The Legend", "The Double Legend", "The Triple Legend", "The Tetra Legend"}
	for _, removeRole := range removeRoles {
		d.Session.GuildMemberRoleRemove()
	}

}
