package ffxiv

type Job struct {
	FullName     string
	Abbreviation string
}

var Jobs = map[string]*Job{
	"Gunbreaker":  {FullName: "Gunbreaker", Abbreviation: "GNB"},
	"Paladin":     {FullName: "Paladin", Abbreviation: "PLD"},
	"Gladiator":   {FullName: "Gladiator", Abbreviation: "GLD"},
	"DarkKnight":  {FullName: "Dark Knight", Abbreviation: "DRK"},
	"Warrior":     {FullName: "Warrior", Abbreviation: "WAR"},
	"Marauder":    {FullName: "Marauder", Abbreviation: "MRD"},
	"Scholar":     {FullName: "Scholar", Abbreviation: "SCH"},
	"Arcanist":    {FullName: "Arcanist", Abbreviation: "ACN"},
	"Sage":        {FullName: "Sage", Abbreviation: "SGE"},
	"Astrologian": {FullName: "Astrologian", Abbreviation: "AST"},
	"WhiteMage":   {FullName: "White Mage", Abbreviation: "WHM"},
	"Conjurer":    {FullName: "Conjurer", Abbreviation: "CNJ"},
	"Samurai":     {FullName: "Samurai", Abbreviation: "SAM"},
	"Dragoon":     {FullName: "Dragoon", Abbreviation: "DRG"},
	"Ninja":       {FullName: "Ninja", Abbreviation: "NIN"},
	"Monk":        {FullName: "Monk", Abbreviation: "MNK"},
	"Reaper":      {FullName: "Reaper", Abbreviation: "RPR"},
	"Bard":        {FullName: "Bard", Abbreviation: "BRD"},
	"Machinist":   {FullName: "Machinist", Abbreviation: "MCH"},
	"Dancer":      {FullName: "Dancer", Abbreviation: "DNC"},
	"BlackMage":   {FullName: "Black Mage", Abbreviation: "BLM"},
	"BlueMage":    {FullName: "Blue Mage", Abbreviation: "BLU"},
	"Summoner":    {FullName: "Summoner", Abbreviation: "SMN"},
	"RedMage":     {FullName: "Red Mage", Abbreviation: "RDM"},
	"Lancer":      {FullName: "Lancer", Abbreviation: "LNC"},
	"Pugilist":    {FullName: "Puligist", Abbreviation: "PUG"},
	"Rogue":       {FullName: "Rogue", Abbreviation: "ROG"},
	"Thaumaturge": {FullName: "Thaumaturge", Abbreviation: "THM"},
	"Archer":      {FullName: "Archer", Abbreviation: "ARC"},
	"Any":         {FullName: "Any", Abbreviation: "ANY"},
}

func (j *Job) IsHealer() bool {
	for _, healer := range []string{"AST", "WHM", "SGE", "SCH"} {
		if j.Abbreviation == healer {
			return true
		}
	}
	return false
}
