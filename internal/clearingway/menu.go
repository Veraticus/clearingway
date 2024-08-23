package clearingway

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// guaranteed UI elements
const(
	MenuMain    string = "menucreate"
	MenuVerify  string = "menuverify"
	MenuRemove  string = "menuremove"
)

// menu types
// currently only default/guaranteed menus, dropdown menu based on encounters
type MenuType string
const (
	MenuTypeDefault   MenuType = "default"
	MenuTypeEncounter MenuType = "encounter"
)

// struct to hold data for all different menu components
type Menus struct {
	Defaults map[string]*Menu
	Menus    []*Menu
}

type Menu struct {
	Name           string  // internal name e.g "menureclear"
	Type           MenuType
	Blurb          string
	Title          string  // title to show in embed
	Description    string  // optional description to show in embed
	ImageURL       string  // optional image URL
	AdditionalData MenuComponent  // additional data depending on MenuType
}

type MenuComponent interface{
	Type() MenuType
}

type MenuRoleHelper struct {
	Role         *Role
	Prerequisite *Role
}

type MenuTypeEncounterData struct {
	Roles        map[string]*MenuRoleHelper
	ExtraRoles   []*Role
	RoleType     []RoleType
	MultiSelect  bool
	RequireClear bool
}

func (MenuTypeEncounterData) Type() MenuType {
	return MenuTypeEncounter
}

func (m *Menu) Init(c *ConfigMenu) {
	m.Name = c.Name
	m.Type = MenuType(c.Type)
	m.Blurb = c.Blurb
	m.Title = c.Title

	if len(c.Description) != 0 {
		m.Description = c.Description
	}

	if len(c.ImageUrl) != 0 {
		m.ImageURL = c.ImageUrl
	}

	switch m.Type {
	case MenuTypeDefault:
	case MenuTypeEncounter:
		m.Type = MenuTypeEncounter
		m.AdditionalData = &MenuTypeEncounterData{}
		data := m.AdditionalData.(*MenuTypeEncounterData)

		if len(c.RoleType) != 0 {
			for _, roleType := range c.RoleType {
				data.RoleType = append(data.RoleType, RoleType(roleType))
			}
		}

		if c.MultiSelect {
			data.MultiSelect = true
		} else {
			data.MultiSelect = false
		}

		if c.RequireClear {
			data.RequireClear = true
		} else {
			data.RequireClear = false
		}

		for _, configRole := range c.ConfigRoles {
			role := &Role{}
			if len(configRole.Name) != 0 {
				role.Name = configRole.Name
			}
			if len(configRole.Description) != 0 {
				role.Description = configRole.Description
			}
			if configRole.Color != 0 {
				role.Color = configRole.Color
			}
			if configRole.Hoist {
				role.Hoist = true
			}
			if configRole.Mention {
				role.Mention = true
			}
			data.ExtraRoles = append(data.ExtraRoles, role)
		}
	}
}

func (g *Guild) DefaultMenus() {
	g.Menus.Defaults[MenuMain] = &Menu{
		Name: MenuMain,
		Type: MenuTypeDefault,
		Title: "Welcome to " + g.Name,
		Description: "Use the buttons below to assign roles!",
	}

	g.Menus.Defaults[MenuVerify] = &Menu{
		Name: MenuVerify,
		Type: MenuTypeDefault,
		Title: "Verify Character",
	}

	g.Menus.Defaults[MenuRemove] = &Menu{
		Name: MenuRemove,
		Type: MenuTypeDefault,
		Title: "Remove Roles",
		Description: "Use the buttons below to remove Clearingway related roles!",
	}
}

// ComponentsHandler is function that executes the respective guild specific functions
func (c *Clearingway) ComponentsHandler(s *discordgo.Session, i *discordgo.InteractionCreate, customID string){
	if g, ok := c.Guilds.Guilds[i.GuildID]; ok {
		g.ComponentsHandlers[customID](s, i)
	} else {
		fmt.Printf("Interaction received from guild %s with no configuration!\n", i.GuildID)
		return
	}
}

// Returns with all the additional roles that were specified
// by the yaml under the menu config
func (ms *Menus) Roles() *Roles {
	roles := &Roles{Roles: []*Role{}}

	for _, menu := range ms.Menus {
		if menu.Type == MenuTypeEncounter {
			roles.Roles = append(roles.Roles, menu.AdditionalData.(*MenuTypeEncounterData).ExtraRoles...)
		}
	}

	return roles
}
