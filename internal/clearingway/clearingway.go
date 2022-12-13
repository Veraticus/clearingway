package clearingway

import (
	"github.com/Veraticus/clearingway/internal/discord"
	"github.com/Veraticus/clearingway/internal/fflogs"
	"github.com/Veraticus/clearingway/internal/ffxiv"

	trie "github.com/Vivino/go-autocomplete-trie"
)

type Clearingway struct {
	Config  *Config
	Discord *discord.Discord
	Guilds  *Guilds
	Fflogs  *fflogs.Fflogs
	Ready   bool

	AllWorlds        []string
	AutoCompleteTrie *trie.Trie
}

func (c *Clearingway) Init() {
	c.AllWorlds = ffxiv.AllWorlds()
	c.AutoCompleteTrie = trie.New()
	for _, world := range c.AllWorlds {
		c.AutoCompleteTrie.Insert(world)
	}

	c.Guilds = &Guilds{Guilds: map[string]*Guild{}}

	for _, configGuild := range c.Config.ConfigGuilds {
		guild := &Guild{}
		guild.Init(configGuild)
		c.Guilds.Guilds[guild.Id] = guild
	}
}
