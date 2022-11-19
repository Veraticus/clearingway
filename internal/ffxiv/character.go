package ffxiv

import (
	"fmt"
	"hash/adler32"
	"strings"
	"time"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Characters struct {
	Characters map[string]*Character
}

type Character struct {
	World          string
	FirstName      string
	LastName       string
	LodestoneID    int
	LastUpdateTime time.Time
}

func (cs *Characters) Init(world, firstName, lastName string) (*Character, error) {
	if len(firstName) < 2 {
		return nil, fmt.Errorf("First name must be at least two characters.")
	}
	if len(lastName) < 2 {
		return nil, fmt.Errorf("Last name must be at least two characters.")
	}
	name := firstName + " " + lastName

	title := cases.Title(language.AmericanEnglish)
	char, ok := cs.Characters[name+"-"+world]
	if !ok {
		char = &Character{
			FirstName: firstName,
			LastName:  lastName,
			World:     title.String(world),
		}

		cs.Characters[name+"-"+world] = char
	}

	return char, nil
}

func (c *Character) UpdatedRecently() bool {
	duration := time.Now().Sub(c.LastUpdateTime)
	return duration.Minutes() <= 5.0
}

func (c *Character) Name() string {
	title := cases.Title(language.AmericanEnglish)
	name := title.String(c.FirstName) + " " + title.String(c.LastName)
	name = strings.Replace(name, "’", "'", 1)
	return name
}

func (c *Character) LodestoneSlug(discordId string) string {
	return fmt.Sprintf("clearingway-%d", adler32.Checksum([]byte(discordId)))
}
