package lodestone

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Veraticus/clearingway/internal/ffxiv"
	"github.com/gocolly/colly"
)

var characterLodestoneUrl = "https://na.finalfantasyxiv.com/lodestone/character/"

func isInList(found_link string, link_list []string) bool {
	for _, stored_link := range link_list {
		if stored_link == found_link {
			return true
		}
	}
	return false
}

func GetAchievements(c *ffxiv.Character) ([]string, error) {
	achievements := []string{}
	visited := []string{}
	links := []string{}
	errors := []error{}

	searchUrl := fmt.Sprintf((characterLodestoneUrl + "%s/achievement"), strconv.Itoa(c.LodestoneID))
	links = append(links, searchUrl)
	scraper := colly.NewCollector()

	for len(links) != 0 {
		scraper.OnError(func(_ *colly.Response, err error) {
			errors = append(errors, err)
		})

		err := scraper.Limit(&colly.LimitRule{
			Delay: 3 * time.Second,
		})
		if err != nil {
			errors = append(errors, err)
		}

		scraper.OnHTML("li.entry", func(e *colly.HTMLElement) {

			r, _ := regexp.Compile("\"(.*?)\"")

			name := r.FindString(string(e.ChildText(".entry__activity__txt")))
			achievement := name

			achievements = append(achievements, achievement)
		})

		scraper.OnHTML("li", func(e *colly.HTMLElement) {
			link := e.ChildAttr("a", "href")
			if !isInList(link, links) && !isInList(link, visited) && strings.Contains(link, characterLodestoneUrl) {
				links = append(links, link)
			}
		})

		visit_link := links[0]
		err = scraper.Visit(visit_link)
		if err != nil {
			errors = append(errors, err)
			return nil, buildError(errors)
		}
		_, links = links[0], links[1:]
		visited = append(visited, visit_link)
	}

	return achievements, nil
}
