package clearingway

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// Strings that will be used to access guild-specific functions
const(
	createOption string = "menucreate"
	verifyOption string = "menuverify"
	reclearOption string = "menureclear"
	progOption string = "menuprog"
	colorOption string = "menucolor"
	removeOption string = "menuremove"
)

// ComponentsHandler is function that executes the respective guild specific functions
func (c *Clearingway) ComponentsHandler(s *discordgo.Session, i *discordgo.InteractionCreate, customID string){
	if g, ok := c.Guilds.Guilds[i.GuildID]; ok {
		g.ComponentsHandlers[customID](s, i)
	} else {
		fmt.Printf("Interaction received from guild %s with no configuration!\n", i.GuildID)
		return
	}
}