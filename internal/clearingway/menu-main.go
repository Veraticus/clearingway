package clearingway

import (
	"fmt"

	"github.com/Veraticus/clearingway/internal/discord"
	"github.com/bwmarrin/discordgo"
)

// Sends the main menu as an standalone message in the channel it is called in
func (c *Clearingway) MenuMainSend(s *discordgo.Session, i *discordgo.InteractionCreate) {
	g, ok := c.Guilds.Guilds[i.GuildID]
	if !ok {
		fmt.Printf("Interaction received from guild %s with no configuration!\n", i.GuildID)
		return
	}

	menu := g.Menus.Menus[string(MenuMain)]
	additionalData := menu.AdditionalData

	_, err := s.ChannelMessageSendComplex(i.ChannelID, additionalData.MessageMainMenu)
	if err != nil {
		fmt.Printf("Error sending Discord message: %v\n", err)
		return
	}

	err = discord.StartInteraction(s, i.Interaction, "Sent menu message.")
	if err != nil {
		fmt.Printf("Error sending Discord message: %v\n", err)
		return
	}
}
