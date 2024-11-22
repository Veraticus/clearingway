package clearingway

import (
	"fmt"
	"strings"

	trie "github.com/Vivino/go-autocomplete-trie"
	"github.com/bwmarrin/discordgo"
)

// types of UI elements
type MenuType string

const (
	MenuMain      MenuType = "menuMain"
	MenuVerify    MenuType = "menuVerify"
	MenuRemove    MenuType = "menuRemove"
	MenuEncounter MenuType = "menuEncounter"
)

type CommandType string

const (
	CommandMenu             CommandType = "menu"
	CommandClearsModal      CommandType = "clearsModal"
	CommandRemoveComfy      CommandType = "removeComfy"
	CommandRemoveColor      CommandType = "removeColor"
	CommandRemoveAll        CommandType = "removeAll"
	CommandRemoveEncounter  CommandType = "removeEncounter"
	CommandEncounterProcess CommandType = "encounterProcess"
)

// struct to hold data for all different menu components
type Menus struct {
	Menus            map[string]*Menu
	Autocomplete     []*discordgo.ApplicationCommandOptionChoice
	AutoCompleteTrie *trie.Trie
}

type Menu struct {
	Name           string                         // internal name to uniquely identify menus
	Type           MenuType                       // type of menu to differentiate AdditionalData types
	Title          string                         // title to show in embed
	Description    string                         // optional description to show in embed
	ImageURL       string                         // optional image URL
	ThumbnailURL   string                         // optional thumbnail URL
	Fields         []*discordgo.MessageEmbedField // embed fields
	AdditionalData *MenuAdditionalData            // additional data depending on MenuType
	Buttons        []discordgo.Button
}

type MenuRoleHelper struct {
	Role         *Role
	Prerequisite *Role
}

type MenuAdditionalData struct {
	MessageMainMenu   *discordgo.MessageSend
	MessageEphemeral  *discordgo.InteractionResponse
	EncounterDropdown []discordgo.SelectMenuOption
	Roles             map[string]*MenuRoleHelper
	ExtraRoles        []*Role
	RoleType          []RoleType
	MultiSelect       bool
	RequireClear      bool
}

func (m *Menu) Init(c *ConfigMenu) {
	m.Name = c.Name
	m.Type = MenuType(c.Type)
	m.Title = c.Title
	m.Fields = []*discordgo.MessageEmbedField{}
	m.Buttons = []discordgo.Button{}

	if len(c.Description) != 0 {
		m.Description = c.Description
	}

	if len(c.ImageUrl) != 0 {
		m.ImageURL = c.ImageUrl
	}

	if len(c.ThumbnailUrl) != 0 {
		m.ThumbnailURL = c.ThumbnailUrl
	}

	if len(c.ConfigFields) != 0 {
		for _, configField := range c.ConfigFields {
			m.Fields = append(m.Fields, &discordgo.MessageEmbedField{
				Name:   configField.Name,
				Value:  configField.Value,
				Inline: configField.Inline,
			})
		}
	}

	switch m.Type {
	case MenuMain:
		for _, configButton := range c.ConfigButtons {
			button := discordgo.Button{}

			if len(configButton.Label) != 0 {
				button.Label = configButton.Label
			} else {
				continue
			}

			if configButton.Style != 0 {
				button.Style = discordgo.ButtonStyle(configButton.Style)
			} else {
				button.Style = discordgo.ButtonStyle(1)
			}

			customIDslice := []string{}
			if len(configButton.MenuName) != 0 && len(configButton.MenuType) != 0 {
				menuType := MenuType(configButton.MenuType)
				switch menuType {
				case MenuVerify:
					fallthrough
				case MenuRemove:
					customIDslice = []string{configButton.MenuType, string(CommandMenu)}
				case MenuEncounter:
					customIDslice = []string{string(MenuEncounter), string(CommandMenu), configButton.MenuName}

				default:
					continue
				}
			} else {
				continue
			}
			button.Disabled = false
			button.CustomID = strings.Join(customIDslice, " ")

			m.Buttons = append(m.Buttons, button)
		}
	case MenuEncounter:
		m.AdditionalData = &MenuAdditionalData{}
		data := m.AdditionalData

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
	g.Menus.Menus[string(MenuMain)] = &Menu{
		Name:        string(MenuMain),
		Type:        MenuMain,
		Title:       "Welcome to " + g.Name,
		Description: "Use the buttons below to assign roles!",
	}

	g.Menus.Menus[string(MenuVerify)] = &Menu{
		Name:  string(MenuVerify),
		Type:  MenuVerify,
		Title: "Verify Character",
	}

	g.Menus.Menus[string(MenuRemove)] = &Menu{
		Name:        string(MenuRemove),
		Type:        MenuRemove,
		Title:       "Remove Roles",
		Description: "Use the buttons below to remove Clearingway related roles!",
	}
}

// Returns with all the additional roles that were specified
// by the yaml under the menu config
func (ms *Menus) Roles() *Roles {
	roles := &Roles{Roles: []*Role{}}

	for _, menu := range ms.Menus {
		if menu.Type == MenuEncounter {
			roles.Roles = append(roles.Roles, menu.AdditionalData.ExtraRoles...)
		}
	}

	return roles
}

func PopulateButtons(buttonsList []discordgo.Button) []discordgo.MessageComponent {
	// count how many buttons to ensure menu doesn't exceed limit of 25 buttons
	ctr := 0
	ret := []discordgo.MessageComponent{}
	for _, button := range buttonsList {
		if ctr >= 25 {
			fmt.Printf("Exceeded button limit, skipping any additional buttons...")
			break
		}
		if ctr%5 == 0 {
			ret = append(ret, discordgo.ActionsRow{Components: []discordgo.MessageComponent{button}})
		} else {
			actionsRow := ret[ctr/5].(discordgo.ActionsRow)
			actionsRow.Components = append(actionsRow.Components, button)
			ret[ctr/5] = actionsRow
		}
		ctr++
	}
	return ret
}

func (m *Menu) FinalizeButtons() {
	components := PopulateButtons(m.Buttons)

	if m.Type == MenuMain {
		m.AdditionalData.MessageMainMenu.Components = components
	} else {
		m.AdditionalData.MessageEphemeral.Data.Components = components
	}
}

// Creates an response to an interaction with a static menu
func (c *Clearingway) MenuStaticRespond(s *discordgo.Session, i *discordgo.InteractionCreate, menuName string) {
	g, ok := c.Guilds.Guilds[i.GuildID]
	if !ok {
		fmt.Printf("Interaction received from guild %s with no configuration!\n", i.GuildID)
		return
	}

	menu := g.Menus.Menus[menuName]
	additionalData := menu.AdditionalData

	err := s.InteractionRespond(i.Interaction, additionalData.MessageEphemeral)
	if err != nil {
		fmt.Printf("Error sending Discord message: %v\n", err)
		return
	}
}

func (c *Clearingway) MenuAutocomplete(s *discordgo.Session, i *discordgo.InteractionCreate) {
	g, ok := c.Guilds.Guilds[i.GuildID]
	if !ok {
		fmt.Printf("Interaction received from guild %s with no configuration!\n", i.GuildID)
		return
	}

	options := i.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	var menu string
	if option, ok := optionMap["menu"]; ok {
		menu = option.StringValue()
	}

	choices := []*discordgo.ApplicationCommandOptionChoice{}

	if len(menu) == 0 {
		choices = g.Menus.Autocomplete
	} else {
		for _, menuCompletion := range g.Menus.AutoCompleteTrie.SearchAll(menu) {
			choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
				Name:  menuCompletion,
				Value: menuCompletion,
			})
		}
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionApplicationCommandAutocompleteResult,
		Data: &discordgo.InteractionResponseData{
			Choices: choices,
		},
	})
	if err != nil {
		fmt.Printf("Could not send Discord autocompletions: %+v\n", err)
	}
}
