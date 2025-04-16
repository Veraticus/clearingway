package clearingway

import (
	"fmt"
	"strings"

	"github.com/Veraticus/clearingway/internal/discord"
	"github.com/bwmarrin/discordgo"
)

func sendMenu(s *discordgo.Session, i *discordgo.InteractionCreate, g *Guild, menuSelection string) error {
	menu, ok := g.Menus.Menus[menuSelection]
	if !ok {
		err := discord.StartInteraction(s, i.Interaction, "Unable to find menu.")
		if err != nil {
			fmt.Printf("Error sending Discord message: %v\n", err)
		}
		return err
	}
	additionalData := menu.AdditionalData

	if menu.Type != MenuMain {
		err := discord.StartInteraction(s, i.Interaction, "Menu is not of type menuMain.")
		if err != nil {
			fmt.Printf("Error sending Discord message: %v\n", err)
		}
		return err
	}

	_, err := s.ChannelMessageSendComplex(i.ChannelID, additionalData.MessageMainMenu)
	if err != nil {
		fmt.Printf("Error sending Discord message: %v\n", err)
		return err
	}

	return nil
}

// Sends the main menu as an standalone message in the channel it is called in
func (c *Clearingway) MenuMainSend(s *discordgo.Session, i *discordgo.InteractionCreate) {
	g, ok := c.Guilds.Guilds[i.GuildID]
	if !ok {
		fmt.Printf("Interaction received from guild %s with no configuration!\n", i.GuildID)
		return
	}

	menuOpt := i.ApplicationCommandData().Options[0].StringValue()

	menuOptSplit := strings.Split(menuOpt, " ")
	if (len(menuOptSplit) > 1) && (menuOptSplit[0] == "group") {
		menuGroup, ok := g.Menus.MenuGroups[menuOptSplit[1]]
		if !ok {
			err := discord.StartInteraction(s, i.Interaction, "Unable to find menu group.")
			if err != nil {
				fmt.Printf("Error sending Discord message: %v\n", err)
			}
			return
		}
		
		for _, menu := range menuGroup {
			err := sendMenu(s, i, g, menu)
			if err != nil {
				fmt.Printf("Error sending menu: %v\n", err)
				continue
			}
		}
	} else {
		err := sendMenu(s, i, g, menuOpt)
		if err != nil {
			fmt.Printf("Error sending menu: %v\n", err)
			return
		}
	}

	err := discord.StartInteraction(s, i.Interaction, "Sent menu message(s).")
	if err != nil {
		fmt.Printf("Error sending Discord message: %v\n", err)
		return
	}
}
