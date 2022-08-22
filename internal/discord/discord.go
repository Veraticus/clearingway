package discord

import (
	"fmt"
	"strings"
	"time"

	"github.com/Veraticus/clearingway/internal/fflogs"
	"github.com/Veraticus/clearingway/internal/ffxiv"

	"github.com/bwmarrin/discordgo"
)

type Discord struct {
	Token     string
	ChannelId string
	Fflogs    *fflogs.Fflogs

	Encounters *fflogs.Encounters
	Roles      *Roles

	Characters *ffxiv.Characters

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

	for _, guild := range event.Guilds {
		existingRoles, err := s.GuildRoles(guild.ID)
		if err != nil {
			fmt.Printf("Error getting existing roles: %v", err)
			return
		}

		fmt.Printf("Ensuring roles...\n")
		err = d.Roles.Ensure(guild.ID, s, existingRoles)
		if err != nil {
			fmt.Printf("Error ensuring roles: %v", err)
			return
		}

		fmt.Printf("Reorder roles...\n")
		err = d.Roles.Reorder(guild.ID, s)
		if err != nil {
			fmt.Printf("Error ensuring roles: %v", err)
			return
		}
	}
	fmt.Printf("Discord ready!\n")
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
		roleNames := d.Roles.RoleNames(m.Member.Roles)
		nick := m.Member.Nick
		if nick == "" {
			nick = m.Author.Username
		}

		char, err := d.Characters.Init(nick, roleNames)
		if err != nil {
			_, err = s.ChannelMessageSendReply(d.ChannelId, fmt.Sprintf("Could not find character: %s", err), (*m).Reference())
			if err != nil {
				fmt.Printf("Error sending Discord message: %v", err)
			}
			return
		}

		message, err := s.ChannelMessageSendReply(d.ChannelId, fmt.Sprintf("Analyzing logs for %s (%s), please wait...", char.Name, char.Server), (*m).Reference())
		if err != nil {
			fmt.Printf("Error sending Discord message: %v", err)
			return
		}

		if char.UpdatedRecently() {
			_, err = s.ChannelMessageEdit(d.ChannelId, message.ID, "Please only use `clears!` once every 5 minutes.")
			if err != nil {
				fmt.Printf("Error sending Discord message: %v", err)
			}
			return
		}

		charText, err := d.UpdateCharacter(char, m.Author.ID, m.GuildID)
		if err != nil {
			_, err = s.ChannelMessageEdit(d.ChannelId, message.ID, fmt.Sprintf("Could not analyze clears for %s (%s): %s", char.Name, char.Server, err))
			if err != nil {
				fmt.Printf("Error sending Discord message: %v", err)
			}
			return
		}

		_, err = s.ChannelMessageEdit(d.ChannelId, message.ID, fmt.Sprintf("Finished clear analysis.\n%s", charText))
		if err != nil {
			fmt.Printf("Error sending Discord message: %v", err)
		}
		return
	}
}

func (d *Discord) UpdateCharacter(char *ffxiv.Character, discordUserId, guildId string) (string, error) {
	encounterRankings, err := d.Fflogs.GetEncounterRankings(d.Encounters, char)
	if err != nil {
		return "", fmt.Errorf("Error retrieving encounter rankings: %w", err)
	}

	member, err := d.Session.GuildMember(guildId, discordUserId)
	if err != nil {
		return "", fmt.Errorf("Could not retrieve roles for user: %w", err)
	}

	text := strings.Builder{}
	text.WriteString("")
	for _, role := range d.Roles.Roles {
		if role.ShouldApply == nil {
			continue
		}

		shouldApply := role.ShouldApply(d.Encounters, encounterRankings)
		if shouldApply {
			if !role.PresentInRoles(member.Roles) {
				err := role.AddToCharacter(guildId, discordUserId, d.Session, char)
				if err != nil {
					return "", fmt.Errorf("Error adding Discord role: %v", err)
				}
				text.WriteString(fmt.Sprintf("Adding role: `%s`\n", role.Name))
			}
		} else {
			if role.PresentInRoles(member.Roles) {
				role.RemoveFromCharacter(guildId, discordUserId, d.Session, char)
				if err != nil {
					return "", fmt.Errorf("Error removing Discord role: %v", err)
				}
				text.WriteString(fmt.Sprintf("Removing role: `%s`\n", role.Name))
			}
		}
	}

	char.LastUpdateTime = time.Now()

	return text.String(), nil
}
