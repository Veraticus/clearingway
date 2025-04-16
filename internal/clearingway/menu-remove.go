package clearingway

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

// Creates the menu component of MenuRemove
func (m *Menu) MenuRemoveInit() {
	removeAllCustomID := []string{string(MenuRemove), string(CommandRemoveAll)}
	removeButtons := []discordgo.Button{
		{
			Label:   "Yes, remove all roles",
			Style: discordgo.DangerButton,
			Disabled: false,
			CustomID: strings.Join(removeAllCustomID, " "),
		},
	}

	m.Buttons = append(m.Buttons, removeButtons...)

	message := &discordgo.InteractionResponseData{
		Flags: discordgo.MessageFlagsEphemeral,
		Embeds: []*discordgo.MessageEmbed{
			{
				Title:       m.Title,
				Description: m.Description,
			},
		},
	}

	m.MenuStyle(message.Embeds)

	m.AdditionalData = &MenuAdditionalData{
		MessageEphemeral: &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: message,
		},
	}
}
