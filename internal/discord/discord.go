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

var (
	clearCommand = &discordgo.ApplicationCommand{
		Name:        "clears",
		Description: "Analyze fflogs and assign yourself cleared roles.",
	}
)

func (d *Discord) Start() error {
	s, err := discordgo.New("Bot " + d.Token)
	if err != nil {
		return fmt.Errorf("Could not start Discord: %f", err)
	}

	s.AddHandler(d.ready)
	s.AddHandler(d.interactionCreate)

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

		fmt.Printf("Adding commands...\n")
		s.ApplicationCommandCreate(event.User.ID, guild.ID, clearCommand)
	}
	fmt.Printf("Discord ready!\n")
}

func (d *Discord) interactionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Ignore messages not on the correct channel
	if i.ChannelID != d.ChannelId {
		return
	}

	// Check if the message is "!clears"
	roleNames := d.Roles.RoleNames(i.Member.Roles)
	nick := i.Member.Nick
	if nick == "" {
		nick = i.Member.User.Username
	}
	char, err := d.Characters.Init(nick, roleNames)

	if err != nil {
		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("Could not find character: %s", err),
			},
		})
		if err != nil {
			fmt.Printf("Error sending Discord message: %v", err)
		}
		return
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("Analyzing logs for %s (%s)...", char.Name, char.Server),
		},
	})
	if err != nil {
		fmt.Printf("Error sending Discord message: %v", err)
		return
	}

	if char.UpdatedRecently() {
		s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Content: "Please only use `/clears` once every 5 minutes.",
		})
		if err != nil {
			fmt.Printf("Error sending Discord message: %v", err)
		}
		return
	}

	charText, err := d.UpdateCharacter(char, i.Member.User.ID, i.GuildID)
	if err != nil {
		s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Content: fmt.Sprintf("Could not analyze clears for %s (%s): %s", char.Name, char.Server, err),
		})
		if err != nil {
			fmt.Printf("Error sending Discord message: %v", err)
		}
		return
	}

	_, err = s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Content: fmt.Sprintf("Finished clear analysis.\n%s", charText),
	})
	if err != nil {
		fmt.Printf("Error sending Discord message: %v", err)
	}
	return
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
