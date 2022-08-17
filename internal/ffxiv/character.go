package ffxiv

import (
	"fmt"
)

type Character struct {
	Name   string
	Server string
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

func InitCharacter(name string, roles []string) (*Character, error) {
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

	character := &Character{
		Name:   name,
		Server: server,
	}

	return character, nil
}
