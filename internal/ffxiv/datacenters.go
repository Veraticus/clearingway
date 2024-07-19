package ffxiv

import (
	"fmt"
)

type PhysicalDatacenter struct {
	Name               string
	Abbreviation       string
	LogicalDatacenters map[string]*LogicalDatacenter
}

type LogicalDatacenter struct {
	Name   string
	Color  int
	Worlds map[string]*World
}

type World struct {
	Name string
}

var NA = &PhysicalDatacenter{
	Name:         "North America",
	Abbreviation: "NA",
	LogicalDatacenters: map[string]*LogicalDatacenter{
		"Aether": {
			Name:  "Aether",
			Color: 0x71368a,
			Worlds: map[string]*World{
				"Adamantoise":  {Name: "Adamantoise"},
				"Cactuar":      {Name: "Cactuar"},
				"Faerie":       {Name: "Faerie"},
				"Gilgamesh":    {Name: "Gilgamesh"},
				"Jenova":       {Name: "Jenova"},
				"Midgardsormr": {Name: "Midgardsormr"},
				"Sargatanas":   {Name: "Sargatanas"},
				"Siren":        {Name: "Siren"},
			},
		},
		"Crystal": {
			Name:  "Crystal",
			Color: 0x206694,
			Worlds: map[string]*World{
				"Balmung":   {Name: "Balmung"},
				"Brynhildr": {Name: "Brynhildr"},
				"Coeurl":    {Name: "Coeurl"},
				"Diabolos":  {Name: "Diabolos"},
				"Goblin":    {Name: "Goblin"},
				"Malboro":   {Name: "Malboro"},
				"Mateus":    {Name: "Mateus"},
				"Zalera":    {Name: "Zalera"},
			},
		},
		"Primal": {
			Name:  "Primal",
			Color: 0x992d22,
			Worlds: map[string]*World{
				"Behemoth":  {Name: "Behemoth"},
				"Excalibur": {Name: "Excalibur"},
				"Exodus":    {Name: "Exodus"},
				"Famfrit":   {Name: "Famfrit"},
				"Hyperion":  {Name: "Hyperion"},
				"Lamia":     {Name: "Lamia"},
				"Leviathan": {Name: "Leviathan"},
				"Ultros":    {Name: "Ultros"},
			},
		},
		"Dynamis": {
			Name:  "Dynamis",
			Color: 0xEAC645,
			Worlds: map[string]*World{
				"Cuchulainn":    {Name: "Cuchulainn"},
				"Golen":         {Name: "Golem"},
				"Halicarnassus": {Name: "Halicarnassus"},
				"Kraken":        {Name: "Kraken"},
				"Maduin":        {Name: "Maduin"},
				"Marilith":      {Name: "Marilith"},
				"Rafflesia":     {Name: "Rafflesia"},
				"Seraph":        {Name: "Seraph"},
			},
		},
	},
}

var EU = &PhysicalDatacenter{
	Name:         "Europe",
	Abbreviation: "EU",
	LogicalDatacenters: map[string]*LogicalDatacenter{
		"Chaos": {
			Name: "Chaos",
			Worlds: map[string]*World{
				"Cerberus":    {Name: "Cerberus"},
				"Louisoix":    {Name: "Louisoix"},
				"Moogle":      {Name: "Moogle"},
				"Omega":       {Name: "Omega"},
				"Phantom":     {Name: "Phantom"},
				"Ragnarok":    {Name: "Ragnarok"},
				"Sagittarius": {Name: "Sagittarius"},
				"Spriggan":    {Name: "Spriggan"},
			},
		},
		"Light": {
			Name: "Light",
			Worlds: map[string]*World{
				"Alpha":     {Name: "Alpha"},
				"Lich":      {Name: "Lich"},
				"Odin":      {Name: "Odin"},
				"Phoenix":   {Name: "Phoenix"},
				"Raiden":    {Name: "Raiden"},
				"Shiva":     {Name: "Shiva"},
				"Twintania": {Name: "Twintania"},
				"Zodiark":   {Name: "Zodiark"},
			},
		},
	},
}

var OC = &PhysicalDatacenter{
	Name:         "Oceania",
	Abbreviation: "OC",
	LogicalDatacenters: map[string]*LogicalDatacenter{
		"Materia": {
			Name: "Materia",
			Worlds: map[string]*World{
				"Bismarck": {Name: "Bismarck"},
				"Ravana":   {Name: "Ravana"},
				"Sephirot": {Name: "Sephirot"},
				"Sophia":   {Name: "Sophia"},
				"Zurvan":   {Name: "Zurvan"},
			},
		},
	},
}

var JP = &PhysicalDatacenter{
	Name:         "Japan",
	Abbreviation: "JP",
	LogicalDatacenters: map[string]*LogicalDatacenter{
		"Elemental": {
			Name: "Elemental",
			Worlds: map[string]*World{
				"Aegis":     {Name: "Aegis"},
				"Atomos":    {Name: "Atomos"},
				"Carbuncle": {Name: "Carbuncle"},
				"Garuda":    {Name: "Garuda"},
				"Gungnir":   {Name: "Gungnir"},
				"Kujata":    {Name: "Kujata"},
				"Tonberry":  {Name: "Tonberry"},
				"Typhon":    {Name: "Typhon"},
			},
		},
		"Gaia": {
			Name: "Gaia",
			Worlds: map[string]*World{
				"Alexander": {Name: "Alexander"},
				"Bahamut":   {Name: "Bahamut"},
				"Durandal":  {Name: "Durandal"},
				"Fenrir":    {Name: "Fenrir"},
				"Ifrit":     {Name: "Ifrit"},
				"Ridill":    {Name: "Ridill"},
				"Tiamat":    {Name: "Tiamat"},
				"Ultima":    {Name: "Ultima"},
			},
		},
		"Mana": {
			Name: "Mana",
			Worlds: map[string]*World{
				"Anima":        {Name: "Anima"},
				"Asura":        {Name: "Asura"},
				"Chocobo":      {Name: "Chocobo"},
				"Hades":        {Name: "Hades"},
				"Ixion":        {Name: "Ixion"},
				"Masamune":     {Name: "Masamune"},
				"Pandaemonium": {Name: "Pandaemonium"},
				"Titan":        {Name: "Titan"},
			},
		},
		"Meteor": {
			Name: "Meteor",
			Worlds: map[string]*World{
				"Belias":     {Name: "Belias"},
				"Mandragora": {Name: "Mandragora"},
				"Ramuh":      {Name: "Ramuh"},
				"Shinryu":    {Name: "Shinryu"},
				"Unicorn":    {Name: "Unicorn"},
				"Valefor":    {Name: "Valefor"},
				"Yojimbo":    {Name: "Yojimbo"},
				"Zeromus":    {Name: "Zeromus"},
			},
		},
	},
}

func WorldsForLogicalDatacenter(logicalDatacenter string) (map[string]*World, error) {
	switch logicalDatacenter {
	case "Aether":
		return NA.LogicalDatacenters["Aether"].Worlds, nil
	case "Primal":
		return NA.LogicalDatacenters["Primal"].Worlds, nil
	case "Crystal":
		return NA.LogicalDatacenters["Crystal"].Worlds, nil
	case "Dynamis":
		return NA.LogicalDatacenters["Dynamis"].Worlds, nil
	case "Chaos":
		return EU.LogicalDatacenters["Chaos"].Worlds, nil
	case "Light":
		return EU.LogicalDatacenters["Light"].Worlds, nil
	case "Materia":
		return OC.LogicalDatacenters["Materia"].Worlds, nil
	case "Elemental":
		return JP.LogicalDatacenters["Elemental"].Worlds, nil
	case "Gaia":
		return JP.LogicalDatacenters["Gaia"].Worlds, nil
	case "Mana":
		return JP.LogicalDatacenters["Mana"].Worlds, nil
	case "Meteor":
		return JP.LogicalDatacenters["Meteor"].Worlds, nil
	default:
		return nil, fmt.Errorf("Could not find datacenter: %v", logicalDatacenter)
	}
}

func AllWorlds() []string {
	worlds := []string{}
	for _, physicalDatacenter := range AllPhysicalDatacenters() {
		for _, logicalDatacenters := range physicalDatacenter.LogicalDatacenters {
			for _, world := range logicalDatacenters.Worlds {
				worlds = append(worlds, world.Name)
			}
		}
	}

	return worlds
}

func IsWorld(w string) bool {
	for _, world := range AllWorlds() {
		if w == world {
			return true
		}
	}
	return false
}

func AllPhysicalDatacenters() []*PhysicalDatacenter {
	return []*PhysicalDatacenter{NA, EU, OC, JP}
}

func PhysicalDatacenterForAbbreviation(name string) *PhysicalDatacenter {
	switch name {
	case "NA":
		return NA
	case "EU":
		return EU
	case "OC":
		return OC
	case "JP":
		return JP
	}
	return nil
}

func PhysicalDatacenterForWorld(name string) *PhysicalDatacenter {
	for _, physicalDataCenter := range AllPhysicalDatacenters() {
		for _, logicalDataCenter := range physicalDataCenter.LogicalDatacenters {
			for _, world := range logicalDataCenter.Worlds {
				if world.Name == name {
					return physicalDataCenter
				}
			}
		}
	}
	return nil
}
