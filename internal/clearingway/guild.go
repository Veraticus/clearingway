package clearingway

import (
	"github.com/Veraticus/clearingway/internal/ffxiv"
)

type Guilds struct {
	Guilds map[string]*Guild
}

type Guild struct {
	Name       string
	Id         string
	ChannelId  string
	Encounters *Encounters
	Characters *ffxiv.Characters

	EncounterRoles *Roles
	ParsingRoles   *Roles
	UltimateRoles  *Roles
	WorldRoles     *Roles
}

func (g *Guild) Init(c *ConfigGuild) {
	g.Name = c.Name
	g.Id = c.GuildId
	g.ChannelId = c.ChannelId
	g.Encounters = &Encounters{Encounters: []*Encounter{}}
	g.Characters = &ffxiv.Characters{Characters: map[string]*ffxiv.Character{}}

	for _, configEncounter := range c.ConfigEncounters {
		encounter := &Encounter{}
		encounter.Init(configEncounter)
		g.Encounters.Encounters = append(g.Encounters.Encounters, encounter)
	}
}

func (g *Guild) AllEncounters() []*Encounter {
	encounters := g.Encounters.Encounters
	encounters = append(encounters, UltimateEncounters.Encounters...)

	return encounters
}

func (g *Guild) NonUltRoles() []*Role {
	roles := g.EncounterRoles.Roles
	roles = append(roles, g.ParsingRoles.Roles...)
	roles = append(roles, g.WorldRoles.Roles...)

	return roles
}
