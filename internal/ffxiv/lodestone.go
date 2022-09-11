package ffxiv

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

var lodestoneUrl = "https://na.finalfantasyxiv.com/lodestone"

func (c *Character) SetLodestoneID() error {
	collector := colly.NewCollector(colly.Async(true))
	collector.SetRequestTimeout(30 * time.Second)
	charIDs := []int{}
	errors := []error{}
	spawnedChildren := false
	searchUrl := fmt.Sprintf(
		"/character/?q=%v&worldname=%v",
		url.QueryEscape(c.Name()),
		c.World,
	)

	collector.OnHTML(".ldst__window .entry", func(e *colly.HTMLElement) {
		name := e.ChildText(".entry__name")
		if strings.ToLower(name) == strings.ToLower(c.Name()) {
			linkText := e.ChildAttr(".entry__link", "href")
			var charID int
			n, err := fmt.Sscanf(linkText, "/lodestone/character/%d/", &charID)
			if n == 0 {
				errors = append(errors, fmt.Errorf("Could not find character ID!"))
			}
			if err != nil {
				errors = append(errors, fmt.Errorf("Could not parse lodestone URL: %w", err))
			}
			charIDs = append(charIDs, charID)
		}
	})

	collector.OnHTML(".ldst__window ul.btn__pager", func(e *colly.HTMLElement) {
		var currentPage int
		var maxPages int
		pages := e.ChildText(".btn__pager__current")
		n, err := fmt.Sscanf(pages, "Page %d of %d", &currentPage, &maxPages)
		if n == 0 {
			errors = append(errors, fmt.Errorf("Could not find pager!"))
		}
		if err != nil {
			errors = append(errors, fmt.Errorf("Could not parse pager: %w", err))
		}
		if !spawnedChildren && currentPage == 1 && maxPages != 1 {
			spawnedChildren = true
			for i := 2; i <= maxPages; i++ {
				e.Request.Visit(lodestoneUrl + searchUrl + fmt.Sprintf("&page=%d", i))
			}
		}
	})

	collector.OnError(func(resp *colly.Response, err error) {
		errors = append(errors, err)
	})

	collector.Visit(lodestoneUrl + searchUrl)
	collector.Wait()

	if len(errors) != 0 {
		return buildError(errors)
	}

	if len(charIDs) == 0 {
		return fmt.Errorf(
			"No character found on the Lodestone for %v (%v)! If you recently renamed yourself or server transferred it can take up to a day for this to be reflected on the Lodestone; please try again later.",
			c.Name(),
			c.World,
		)
	}
	if len(charIDs) > 1 {
		return fmt.Errorf(
			"Too many characters found for name %v! Ensure it is exactly your character name.",
			c.Name(),
		)
	}

	c.LodestoneID = charIDs[0]

	return nil
}

func (c *Character) IsOwner(discordId string) (bool, error) {
	collector := colly.NewCollector(colly.Async(true))
	collector.SetRequestTimeout(30 * time.Second)
	errors := []error{}
	bio := ""

	collector.OnHTML(".character__content.selected", func(e *colly.HTMLElement) {
		bio = e.ChildText(".character__selfintroduction")
	})

	collector.OnError(func(resp *colly.Response, err error) {
		errors = append(errors, err)
	})

	collector.Visit(lodestoneUrl + fmt.Sprintf("/character/%d/", c.LodestoneID))
	collector.Wait()

	if len(errors) != 0 {
		return false, buildError(errors)
	}

	if !strings.Contains(bio, c.LodestoneSlug(discordId)) {
		return false, nil
	}

	// "Your Lodestone profile can be edited from: https://na.finalfantasyxiv.com/lodestone/my/setting/profile/"
	return true, nil
}

func buildError(errors []error) error {
	errorText := strings.Builder{}
	for _, e := range errors {
		errorText.WriteString(e.Error() + "\n")
	}
	return fmt.Errorf("Encountered search errors:\n%v", errorText.String())
}
