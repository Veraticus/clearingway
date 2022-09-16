package ffxiv

var AetherWorlds = map[string]interface{}{
	"Adamantoise":  nil,
	"Cactuar":      nil,
	"Faerie":       nil,
	"Gilgamesh":    nil,
	"Jenova":       nil,
	"Midgardsormr": nil,
	"Sargatanas":   nil,
	"Siren":        nil,
}

var CrystalWorlds = map[string]interface{}{
	"Balmung":   nil,
	"Brynhildr": nil,
	"Coeurl":    nil,
	"Diabolos":  nil,
	"Goblin":    nil,
	"Malboro":   nil,
	"Mateus":    nil,
	"Zalera":    nil,
}

var PrimalWorlds = map[string]interface{}{
	"Behemoth":  nil,
	"Excalibur": nil,
	"Exodus":    nil,
	"Famfrit":   nil,
	"Hyperion":  nil,
	"Lamia":     nil,
	"Leviathan": nil,
	"Ultros":    nil,
}

func AllWorlds() []string {
	worlds := []string{}
	for world := range AetherWorlds {
		worlds = append(worlds, world)
	}
	for world := range CrystalWorlds {
		worlds = append(worlds, world)
	}
	for world := range PrimalWorlds {
		worlds = append(worlds, world)
	}
	return worlds
}
