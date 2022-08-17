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
	Token              string
	ChannelId          string
	Fflogs             *fflogs.Fflogs
	RelevantEncounters *fflogs.Encounters

	Roles *Roles

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
		char, err := d.Characters.Init(m.Member.Nick, roleNames)
		if err != nil {
			_, err = s.ChannelMessageSendReply(d.ChannelId, fmt.Sprintf("Could not find character %s: %s", m.Member.Nick, err), (*m).Reference())
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
			_, err = s.ChannelMessageEdit(d.ChannelId, message.ID, "Please only use `cleared!` once every 5 minutes.")
			if err != nil {
				fmt.Printf("Error sending Discord message: %v", err)
			}
			return
		}

		encounterIds := []int{}
		encounterIds = append(encounterIds, d.RelevantEncounters.IDs()...)
		encounterIds = append(encounterIds, fflogs.UltimateEncounters.IDs()...)
		encounterRankings, err := d.Fflogs.GetEncounterRankings(encounterIds, char)
		if err != nil {
			_, err = s.ChannelMessageEdit(d.ChannelId, message.ID, fmt.Sprintf("Error retrieving encounter rankings: %s", err))
			if err != nil {
				fmt.Printf("Error sending Discord message: %v", err)
			}
			return
		}

		text := strings.Builder{}
		text.WriteString("")
		for _, role := range d.Roles.Roles {
			if role.ShouldApply == nil {
				continue
			}

			shouldApply := role.ShouldApply(d.RelevantEncounters, encounterRankings)
			if shouldApply {
				if !role.PresentInRoles(m.Member.Roles) {
					err := role.AddToCharacter(m.GuildID, m.Author.ID, s, char)
					if err != nil {
						fmt.Printf("Error adding Discord role: %v", err)
					}
					text.WriteString(fmt.Sprintf("Adding role: `%s`\n", role.Name))
				}
			} else {
				if role.PresentInRoles(m.Member.Roles) {
					role.RemoveFromCharacter(m.GuildID, m.Author.ID, s, char)
					if err != nil {
						fmt.Printf("Error removing Discord role: %v", err)
					}
					text.WriteString(fmt.Sprintf("Removing role: `%s`\n", role.Name))
				}
			}
		}

		char.LastUpdateTime = time.Now()
		_, err = s.ChannelMessageEdit(d.ChannelId, message.ID, fmt.Sprintf("Finished clear analysis.\n%s", text.String()))
		if err != nil {
			fmt.Printf("Error sending Discord message: %v", err)
		}
		return
	}
}
