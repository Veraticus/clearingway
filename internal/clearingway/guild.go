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

	RelevantParsingEnabled bool
	RelevantFlexingEnabled bool
	LegendEnabled          bool
	UltimateFlexingEnabled bool
	WorldEnabled           bool

	EncounterRoles       *Roles
	RelevantParsingRoles *Roles
	RelevantFlexingRoles *Roles
	LegendRoles          *Roles
	UltimateFlexingRoles *Roles
	WorldRoles           *Roles
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

	if c.ConfigRoles != nil && c.ConfigRoles.RelevantParsing == false {
		g.RelevantParsingEnabled = false
	} else {
		g.RelevantParsingEnabled = true
	}

	if c.ConfigRoles != nil && c.ConfigRoles.RelevantFlexing == false {
		g.RelevantFlexingEnabled = false
	} else {
		g.RelevantFlexingEnabled = true
	}

	if c.ConfigRoles != nil && c.ConfigRoles.Legend == false {
		g.LegendEnabled = false
	} else {
		g.LegendEnabled = true
	}

	if c.ConfigRoles != nil && c.ConfigRoles.UltimateFlexing == false {
		g.UltimateFlexingEnabled = false
	} else {
		g.UltimateFlexingEnabled = true
	}

	if c.ConfigRoles != nil && c.ConfigRoles.World == false {
		g.WorldEnabled = false
	} else {
		g.WorldEnabled = true
	}
}

func (g *Guild) AllEncounters() []*Encounter {
	encounters := g.Encounters.Encounters
	encounters = append(encounters, UltimateEncounters.Encounters...)

	return encounters
}

func (g *Guild) NonUltRoles() []*Role {
	roles := g.EncounterRoles.Roles

	if g.RelevantParsingEnabled {
		roles = append(roles, g.RelevantParsingRoles.Roles...)
	}
	if g.RelevantFlexingEnabled {
		roles = append(roles, g.RelevantFlexingRoles.Roles...)
	}
	if g.WorldEnabled {
		roles = append(roles, g.WorldRoles.Roles...)
	}

	return roles
}

func (g *Guild) UltRoles() []*Role {
	roles := []*Role{}

	if g.LegendEnabled {
		roles = append(roles, g.LegendRoles.Roles...)
	}
	if g.UltimateFlexingEnabled {
		roles = append(roles, g.UltimateFlexingRoles.Roles...)
	}

	return roles
}
