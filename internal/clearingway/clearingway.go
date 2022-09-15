package clearingway

import (
	"github.com/Veraticus/clearingway/internal/discord"
	"github.com/Veraticus/clearingway/internal/fflogs"
)

type Clearingway struct {
	Config  *Config
	Discord *discord.Discord
	Guilds  *Guilds
	Fflogs  *fflogs.Fflogs
}

func (c *Clearingway) Init() {
	c.Guilds = &Guilds{Guilds: map[string]*Guild{}}

	for _, configGuild := range c.Config.ConfigGuilds {
		guild := &Guild{}
		guild.Init(configGuild)
		c.Guilds.Guilds[guild.Id] = guild
	}
}
