package discord

import (
	"fmt"
	"strings"
	"time"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

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

var verifyCommand = &discordgo.ApplicationCommand{
	Name:        "verify",
	Description: "Verify you own your character and assign them roles.",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "world",
			Description: "Your character's world",
			Required:    true,
		},
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "first-name",
			Description: "Your character's first name",
			Required:    true,
		},
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "last-name",
			Description: "Your character's last name",
			Required:    true,
		},
	},
}

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
			fmt.Printf("Error getting existing roles: %v\n", err)
			return
		}

		fmt.Printf("Ensuring roles...\n")
		err = d.Roles.Ensure(guild.ID, s, existingRoles)
		if err != nil {
			fmt.Printf("Error ensuring roles: %v\n", err)
			return
		}

		fmt.Printf("Reorder roles...\n")
		err = d.Roles.Reorder(guild.ID, s)
		if err != nil {
			fmt.Printf("Error reordering roles: %v\n", err)
			return
		}

		fmt.Printf("Adding commands...\n")
		_, err = s.ApplicationCommandCreate(event.User.ID, guild.ID, verifyCommand)
		if err != nil {
			fmt.Printf("Could not add command: %v\n", err)
		}
	}
	fmt.Printf("Discord ready!\n")
}

func (d *Discord) interactionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Ignore messages not on the correct channel
	if i.ChannelID != d.ChannelId {
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Only issue this command in #botspam!",
			},
		})
		if err != nil {
			fmt.Printf("Error sending Discord message: %v", err)
			return
		}
	}

	options := i.ApplicationCommandData().Options

	// Or convert the slice into a map
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	var world string
	var firstName string
	var lastName string

	if option, ok := optionMap["world"]; ok {
		world = option.StringValue()
	}
	if option, ok := optionMap["first-name"]; ok {
		firstName = option.StringValue()
	}
	if option, ok := optionMap["last-name"]; ok {
		lastName = option.StringValue()
	}

	title := cases.Title(language.AmericanEnglish)
	world = title.String(world)
	firstName = title.String(firstName)
	lastName = title.String(lastName)

	if len(world) == 0 || len(firstName) == 0 || len(lastName) == 0 {
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "/verify command failed! Make sure you input your world, first name, and last name.",
			},
		})
		if err != nil {
			fmt.Printf("Error sending Discord message: %v", err)
		}
		return
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("Finding %s %s (%s) in the Lodestone...", firstName, lastName, world),
		},
	})
	if err != nil {
		fmt.Printf("Error sending Discord message: %v", err)
		return
	}

	char, err := d.Characters.Init(world, firstName, lastName)
	if err != nil {
		_, err = s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Content: err.Error(),
		})
		if err != nil {
			fmt.Printf("Error sending Discord message: %v", err)
		}
		return
	}

	_, err = s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Content: fmt.Sprintf("Verifying ownership of %s (%s)...", char.Name(), char.World),
	})
	if err != nil {
		fmt.Printf("Error sending Discord message: %v", err)
	}

	discordId := i.Member.User.ID
	isOwner, err := char.IsOwner(discordId)
	if err != nil {
		_, err = s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Content: err.Error(),
		})
		if err != nil {
			fmt.Printf("Error sending Discord message: %v", err)
		}
		return
	}
	if !isOwner {
		_, err = s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Content: fmt.Sprintf(
				"You are not the owner of %s (%s)!\nIf this is your character, add the following code to your Lodestone profile:\n\n**%s**\n\nYou can edit your Lodestone profile at https://na.finalfantasyxiv.com/lodestone/my/setting/profile/",
				char.Name(),
				char.World,
				char.LodestoneSlug(discordId),
			),
		})
		if err != nil {
			fmt.Printf("Error sending Discord message: %v", err)
		}
		return
	}

	_, err = s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Content: fmt.Sprintf("Analyzing logs for %s (%s)...", char.Name(), char.World),
	})
	if err != nil {
		fmt.Printf("Error sending Discord message: %v", err)
	}

	if char.UpdatedRecently() {
		_, err = s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Content: "Finished clear analysis.",
		})
		if err != nil {
			fmt.Printf("Error sending Discord message: %v", err)
		}
		return
	}

	charText, err := d.UpdateCharacter(char, i.Member.User.ID, i.GuildID)
	if err != nil {
		_, err = s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Content: fmt.Sprintf("Could not analyze clears for %s (%s): %s", char.Name(), char.World, err),
		})
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
