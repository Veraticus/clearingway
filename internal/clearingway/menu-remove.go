package clearingway

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

// Creates the menu component of MenuRemove
func (m *Menu) MenuRemoveInit() {
	type removeButton struct {
		name        string
		commandType CommandType
	}

	removeButtons := []discordgo.MessageComponent{}
	removeButtonsList := []removeButton{
		{name: "Uncomfy", commandType: CommandRemoveComfy},
		{name: "Uncolor", commandType: CommandRemoveColor},
		{name: "Remove All", commandType: CommandRemoveAll},
	}

	for _, button := range removeButtonsList {
		customIDslice := []string{string(MenuRemove), string(button.commandType)}
		removeButtons = append(removeButtons, discordgo.Button{
			Label:    button.name,
			Style:    discordgo.DangerButton,
			Disabled: false,
			CustomID: strings.Join(customIDslice, " "),
		})
	}

	message := &discordgo.InteractionResponseData{
		Flags: discordgo.MessageFlagsEphemeral,
		Embeds: []*discordgo.MessageEmbed{
			{
				Title:       m.Title,
				Description: m.Description,
			},
		},
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: removeButtons,
			},
		},
	}

	if len(m.ImageURL) > 0 {
		message.Embeds[0].Image = &discordgo.MessageEmbedImage{URL: m.ImageURL}
	}

	m.AdditionalData = &MenuAdditionalData{
		MessageEphemeral: &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: message,
		},
	}

}

func (m *Menu) MenuRemoveAddButton(menuEncounterToAdd *Menu) {
	// create button
	customIDslice := []string{string(MenuRemove), string(CommandRemoveEncounter), string(menuEncounterToAdd.Name)}
	button := discordgo.Button{
		Label:    "Remove " + menuEncounterToAdd.Title,
		Style:    discordgo.DangerButton,
		Disabled: false,
		CustomID: strings.Join(customIDslice, " "),
	}

	// create new ActionsRow component if already made, append to it
	buttonRows := &m.AdditionalData.MessageEphemeral.Data.Components
	if len(*buttonRows) < 2 {
		*buttonRows = append(*buttonRows, discordgo.ActionsRow{Components: []discordgo.MessageComponent{button}})
	} else {
		// creates a copy of ActionsRow, appends the button, then sets the element to the appended copy
		addButtonRow := (*buttonRows)[1].(discordgo.ActionsRow)
		addButtonRow.Components = append(addButtonRow.Components, button)
		(*buttonRows)[1] = addButtonRow
	}
}
