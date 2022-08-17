package ffxiv

import (
	"fmt"
	"time"
)

type Characters struct {
	Characters map[string]*Character
}

type Character struct {
	Name           string
	Server         string
	LastUpdateTime time.Time
}

var AetherServers = []string{
	"Adamantoise",
	"Cactuar",
	"Faerie",
	"Gilgamesh",
	"Jenova",
	"Midgardsormr",
	"Sargatanas",
	"Siren",
}

var CrystalServers = []string{
	"Balmung",
	"Brynhildr",
	"Coeurl",
	"Diabolos",
	"Goblin",
	"Malboro",
	"Mateus",
	"Zalera",
}

var PrimalServers = []string{
	"Behemoth",
	"Excalibur",
	"Exodus",
	"Famfrit",
	"Hyperion",
	"Lamia",
	"Leviathan",
	"Ultros",
}

var NAServers = append(append(append([]string{}, AetherServers...), CrystalServers...), PrimalServers...)

func (cs *Characters) Init(name string, roles []string) (*Character, error) {
	server := ""
	for _, role := range roles {
		for _, naServer := range NAServers {
			if role == naServer {
				server = role
				break
			}
		}
	}
	if server == "" {
		return nil, fmt.Errorf("No NA server role found! Did you run `!iam verify`?")
	}

	char, ok := cs.Characters[name+"-"+server]
	if !ok {
		char = &Character{
			Name:   name,
			Server: server,
		}
		cs.Characters[name+"-"+server] = char
	}

	return char, nil
}

func (c *Character) UpdatedRecently() bool {
	duration := time.Now().Sub(c.LastUpdateTime)
	return duration.Minutes() <= 5.0
}
