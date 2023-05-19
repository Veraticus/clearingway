package clearingway

import (
	"fmt"

	"github.com/Veraticus/clearingway/internal/ffxiv"
)

type Guilds struct {
	Guilds map[string]*Guild
}

type Guild struct {
	Name                string
	Id                  string
	ChannelId           string
	Encounters          *Encounters
	Characters          *ffxiv.Characters
	PhysicalDatacenters *PhysicalDatacenters

	RelevantParsingEnabled    bool
	RelevantFlexingEnabled    bool
	RelevantRepetitionEnabled bool
	LegendEnabled             bool
	UltimateFlexingEnabled    bool
	UltimateRepetitionEnabled bool
	DatacenterEnabled         bool
	SkipRemoval               bool

	EncounterRoles          *Roles
	RelevantParsingRoles    *Roles
	RelevantFlexingRoles    *Roles
	RelevantRepetitionRoles *Roles
	LegendRoles             *Roles
	UltimateFlexingRoles    *Roles
	UltimateRepetitionRoles *Roles
	DatacenterRoles         *Roles
}

func (g *Guild) Init(c *ConfigGuild) {
	g.Name = c.Name
	g.Id = c.GuildId
	g.ChannelId = c.ChannelId
	g.Encounters = &Encounters{Encounters: []*Encounter{}}
	g.Characters = &ffxiv.Characters{Characters: map[string]*ffxiv.Character{}}

	g.PhysicalDatacenters = &PhysicalDatacenters{PhysicalDatacenters: map[string]*PhysicalDatacenter{}}
	fmt.Printf("Datacenters are %+v\n", c.ConfigPhysicalDatacenters)
	g.PhysicalDatacenters.Init(c.ConfigPhysicalDatacenters)

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

	if c.ConfigRoles != nil && c.ConfigRoles.RelevantRepetition == false {
		g.RelevantRepetitionEnabled = false
	} else {
		g.RelevantRepetitionEnabled = true
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

	if c.ConfigRoles != nil && c.ConfigRoles.UltimateRepetition == false {
		g.UltimateRepetitionEnabled = false
	} else {
		g.UltimateRepetitionEnabled = true
	}

	if c.ConfigRoles != nil && c.ConfigRoles.Datacenter == false {
		g.DatacenterEnabled = false
	} else {
		g.DatacenterEnabled = true
	}

	if c.ConfigRoles != nil && c.ConfigRoles.SkipRemoval == true {
		g.SkipRemoval = true
	} else {
		g.SkipRemoval = false
	}

	g.EncounterRoles = g.Encounters.Roles()

	if g.RelevantParsingEnabled {
		g.RelevantParsingRoles = RelevantParsingRoles()
	}

	if g.RelevantFlexingEnabled {
		g.RelevantFlexingRoles = RelevantFlexingRoles()
	}

	if g.RelevantRepetitionEnabled {
		g.RelevantRepetitionRoles = RelevantRepetitionRoles(g.Encounters)
	}

	if g.LegendEnabled {
		g.LegendRoles = LegendRoles()
	}

	if g.UltimateFlexingEnabled {
		g.UltimateFlexingRoles = UltimateFlexingRoles()
	}

	if g.UltimateRepetitionEnabled {
		g.UltimateRepetitionRoles = UltimateRepetitionRoles()
	}

	if g.DatacenterEnabled {
		g.DatacenterRoles = g.PhysicalDatacenters.AllRoles()
	}

	if len(c.ConfigReconfigureRoles) != 0 {
		for _, configReconfigureRole := range c.ConfigReconfigureRoles {
			for _, role := range g.AllRoles() {
				if role.Name == configReconfigureRole.From {
					if configReconfigureRole.To != "" {
						role.Name = configReconfigureRole.To
					}
					if configReconfigureRole.Color != 0 {
						role.Color = configReconfigureRole.Color
					}
					if configReconfigureRole.Skip == true {
						role.Skip = true
					}
					if configReconfigureRole.DontSkip == true {
						role.Skip = false
					}
				}
			}
		}
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
	if g.RelevantRepetitionEnabled {
		roles = append(roles, g.RelevantRepetitionRoles.Roles...)
	}
	if g.DatacenterEnabled {
		roles = append(roles, g.DatacenterRoles.Roles...)
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
	if g.UltimateRepetitionEnabled {
		roles = append(roles, g.UltimateRepetitionRoles.Roles...)
	}

	return roles
}

func (g *Guild) AllRoles() []*Role {
	return append(g.NonUltRoles(), g.UltRoles()...)
}

func (g *Guild) IsProgEnabled() bool {
	for _, encounter := range g.Encounters.Encounters {
		if encounter.ProgRoles != nil {
			return true
		}
	}

	return false
}
