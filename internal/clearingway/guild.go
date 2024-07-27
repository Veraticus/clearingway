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
	Achievements        *Achievements
	Characters          *ffxiv.Characters
	PhysicalDatacenters *PhysicalDatacenters

	RelevantParsingEnabled    bool
	RelevantFlexingEnabled    bool
	RelevantRepetitionEnabled bool
	LegendEnabled             bool
	UltimateFlexingEnabled    bool
	UltimateRepetitionEnabled bool
	DatacenterEnabled         bool
	NameColorsEnabled         bool
	ReclearsEnabled	          bool
	SkipRemoval               bool

	EncounterRoles          *Roles
	RelevantParsingRoles    *Roles
	RelevantFlexingRoles    *Roles
	RelevantRepetitionRoles *Roles
	LegendRoles             *Roles
	UltimateFlexingRoles    *Roles
	UltimateRepetitionRoles *Roles
	DatacenterRoles         *Roles
	AchievementRoles        *Roles
}

func (g *Guild) Init(c *ConfigGuild) {
	g.Name = c.Name
	g.Id = c.GuildId
	g.ChannelId = c.ChannelId
	g.Encounters = &Encounters{Encounters: []*Encounter{}}
	g.Achievements = &Achievements{Achievements: []*Achievement{}}
	g.Characters = &ffxiv.Characters{Characters: map[string]*ffxiv.Character{}}

	g.PhysicalDatacenters = &PhysicalDatacenters{PhysicalDatacenters: map[string]*PhysicalDatacenter{}}
	fmt.Printf("Datacenters are %+v\n", c.ConfigPhysicalDatacenters)
	g.PhysicalDatacenters.Init(c.ConfigPhysicalDatacenters)

	for _, configEncounter := range c.ConfigEncounters {
		encounter := &Encounter{}
		encounter.Init(configEncounter)
		g.Encounters.Encounters = append(g.Encounters.Encounters, encounter)
	}

	for _, configAchievement := range c.ConfigAchievements {
		achievement := &Achievement{}
		achievement.Init(configAchievement)
		g.Achievements.Achievements = append(g.Achievements.Achievements, achievement)
	}

	if c.ConfigRoles != nil {
		if !c.ConfigRoles.RelevantParsing {
			g.RelevantParsingEnabled = false
		} else {
			g.RelevantParsingEnabled = true
		}

		if !c.ConfigRoles.RelevantFlexing {
			g.RelevantFlexingEnabled = false
		} else {
			g.RelevantFlexingEnabled = true
		}

		if !c.ConfigRoles.RelevantRepetition {
			g.RelevantRepetitionEnabled = false
		} else {
			g.RelevantRepetitionEnabled = true
		}

		if !c.ConfigRoles.Legend {
			g.LegendEnabled = false
		} else {
			g.LegendEnabled = true
		}

		if !c.ConfigRoles.UltimateFlexing {
			g.UltimateFlexingEnabled = false
		} else {
			g.UltimateFlexingEnabled = true
		}

		if !c.ConfigRoles.UltimateRepetition {
			g.UltimateRepetitionEnabled = false
		} else {
			g.UltimateRepetitionEnabled = true
		}

		if !c.ConfigRoles.Datacenter {
			g.DatacenterEnabled = false
		} else {
			g.DatacenterEnabled = true
		}

		if c.ConfigRoles.NameColor {
			g.NameColorsEnabled = true
		} else {
			g.NameColorsEnabled = false
		}
		
		if c.ConfigRoles.Reclear {
			g.ReclearsEnabled = true
		} else {
			g.ReclearsEnabled = false
		}

		if c.ConfigRoles.SkipRemoval {
			g.SkipRemoval = true
		} else {
			g.SkipRemoval = false
		}
	}

	g.EncounterRoles = g.Encounters.Roles()
	g.AchievementRoles = g.Achievements.Roles()

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
					// If additional constraints are on this reconfigureRole,
					// make sure they match
					if configReconfigureRole.Type != "" && string(role.Type) != configReconfigureRole.Type {
						continue
					}

					if configReconfigureRole.EncounterName != "" && role.Encounter.Name != configReconfigureRole.EncounterName {
						continue
					}

					if configReconfigureRole.To != "" {
						role.Name = configReconfigureRole.To
					}
					if configReconfigureRole.Color != 0 {
						role.Color = configReconfigureRole.Color
					}
					if configReconfigureRole.Skip {
						role.Skip = true
					}
					if configReconfigureRole.DontSkip {
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
	roles = append(roles, g.AchievementRoles.Roles...)

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
