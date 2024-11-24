package clearingway

import (
	"fmt"
	"strings"

	"github.com/Veraticus/clearingway/internal/discord"
	"github.com/bwmarrin/discordgo"
)

// MenuVerifySendModal sends the user a modal that asks for their character's
// first name, last name, and world to verify their clears
func MenuVerifySendModal(s *discordgo.Session, i *discordgo.InteractionCreate) {
	customID := []string{string(MenuVerify), string(CommandClearsModal)}
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID: strings.Join(customID, " "),
			Title:    "Verify your clears",
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:    "firstName",
							Style:       discordgo.TextInputShort,
							Label:       "Character first name",
							Placeholder: "First name",
							Required:    true,
							MaxLength:   15,
							MinLength:   2,
						},
					},
				},
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:    "lastName",
							Style:       discordgo.TextInputShort,
							Label:       "Character last name",
							Placeholder: "Last name",
							Required:    true,
							MaxLength:   15,
							MinLength:   2,
						},
					},
				},
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:    "world",
							Style:       discordgo.TextInputShort,
							Label:       "Character world",
							Placeholder: "World",
							Required:    true,
							MaxLength:   20,
							MinLength:   2,
						},
					},
				},
			},
		},
	})
	if err != nil {
		fmt.Printf("Unable to send modal: %v", err)
		return
	}
}

// MenuVerifyProcess takes the submitted data from the modal above and processes the character
func (c *Clearingway) MenuVerifyProcess(s *discordgo.Session, i *discordgo.InteractionCreate) {
	g, ok := c.Guilds.Guilds[i.GuildID]
	if !ok {
		fmt.Printf("Interaction received from guild %s with no configuration!\n", i.GuildID)
		return
	}

	// Get the text fields
	options := i.ModalSubmitData()

	err := discord.StartInteraction(s, i.Interaction, "Received character info...")
	if err != nil {
		fmt.Printf("Error sending Discord message: %v\n", err)
		return
	}

	var world string
	var firstName string
	var lastName string

	firstName = options.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
	lastName = options.Components[1].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
	world = options.Components[2].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value

	c.ClearsHelper(s, i, g, world, firstName, lastName)
}
