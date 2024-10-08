package clearingway

import (
	"fmt"
	"strings"

	"github.com/Veraticus/clearingway/internal/discord"
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
	CommandEncounterProcess CommandType = "encounterProcess"
)

// struct to hold data for all different menu components
type Menus struct {
	Menus map[string]*Menu
}

type Menu struct {
	Name           string              // internal name to uniquely identify menus
	Type           MenuType            // type of menu to differentiate AdditionalData types
	Title          string              // title to show in embed
	Description    string              // optional description to show in embed
	ImageURL       string              // optional image URL
	AdditionalData *MenuAdditionalData // additional data depending on MenuType
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

	if len(c.Description) != 0 {
		m.Description = c.Description
	}

	if len(c.ImageUrl) != 0 {
		m.ImageURL = c.ImageUrl
	}

	switch m.Type {
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

// Creates the menu component of MenuRemove
func (m *Menu) MenuRemoveInit() {
	type removeButton struct {
		name        string
		commandType CommandType
	}

	removeButtons := []discordgo.MessageComponent{}
	removeButtonsList := []removeButton{
		{name: "Uncomfy", commandType: CommandRemoveComfy},
		{name: "Uncolor", commandType: CommandRemoveColor},
		{name: "Remove All", commandType: CommandRemoveAll},
	}

	for _, button := range removeButtonsList {
		customIDslice := []string{string(MenuRemove), string(button.commandType)}
		removeButtons = append(removeButtons, discordgo.Button{
			Label:    button.name,
			Style:    discordgo.DangerButton,
			Disabled: false,
			CustomID: strings.Join(customIDslice, " "),
		})
	}

	message := &discordgo.InteractionResponseData{
		Flags: discordgo.MessageFlagsEphemeral,
		Embeds: []*discordgo.MessageEmbed{
			{
				Title:       m.Title,
				Description: m.Description,
			},
		},
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: removeButtons,
			},
		},
	}

	if len(m.ImageURL) > 0 {
		message.Embeds[0].Image = &discordgo.MessageEmbedImage{URL: m.ImageURL}
	}

	m.AdditionalData = &MenuAdditionalData{
		MessageEphemeral: &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: message,
		},
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

// Sends the main menu as an standalone message in the channel it is called in
func (c *Clearingway) MenuMainSend(s *discordgo.Session, i *discordgo.InteractionCreate) {
	g, ok := c.Guilds.Guilds[i.GuildID]
	if !ok {
		fmt.Printf("Interaction received from guild %s with no configuration!\n", i.GuildID)
		return
	}

	menu := g.Menus.Menus[string(MenuMain)]
	additionalData := menu.AdditionalData

	_, err := s.ChannelMessageSendComplex(i.ChannelID, additionalData.MessageMainMenu)
	if err != nil {
		fmt.Printf("Error sending Discord message: %v\n", err)
		return
	}

	err = discord.StartInteraction(s, i.Interaction, "Sent menu message.")
	if err != nil {
		fmt.Printf("Error sending Discord message: %v\n", err)
		return
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

func (m *Menu) MenuEncounterInit(es *Encounters, roleTypes []RoleType) {
	additionalData := m.AdditionalData
	additionalData.Roles = make(map[string]*MenuRoleHelper)
	dropdownSlice := []discordgo.SelectMenuOption{}

	// generate options based on encounters
	for _, roleType := range roleTypes {
		for _, encounter := range es.Encounters {
			role, ok := encounter.Roles[roleType]
			if !ok {
				fmt.Printf("Menu %v: role type %v for encounter %v not found. Skipping...\n", m.Name, string(roleType), encounter.Name)
				continue
			}

			dropdownOption := discordgo.SelectMenuOption{
				Label:   role.Name,
				Value:   role.DiscordRole.ID,
				Default: false,
			}

			menuRoleHelper := &MenuRoleHelper{
				Role: role,
			}

			if additionalData.RequireClear {
				if prereq, ok := encounter.Roles[ClearedRole]; ok {
					menuRoleHelper.Prerequisite = prereq
				}
			}

			dropdownSlice = append(dropdownSlice, dropdownOption)
			additionalData.Roles[role.DiscordRole.ID] = menuRoleHelper
		}
	}

	// generate options based on additional roles
	for _, role := range additionalData.ExtraRoles {
		dropdownOption := discordgo.SelectMenuOption{
			Label:   role.Name,
			Value:   role.DiscordRole.ID,
			Default: false,
		}

		menuRoleHelper := &MenuRoleHelper{
			Role: role,
		}

		dropdownSlice = append(dropdownSlice, dropdownOption)
		additionalData.Roles[role.DiscordRole.ID] = menuRoleHelper
	}

	additionalData.EncounterDropdown = dropdownSlice
}

func (c *Clearingway) MenuEncounterSend(s *discordgo.Session, i *discordgo.InteractionCreate, menuName string) {
	g, ok := c.Guilds.Guilds[i.GuildID]
	if !ok {
		fmt.Printf("Interaction received from guild %s with no configuration!\n", i.GuildID)
		return
	}

	userRoles := i.Member.Roles
	userRoleMap := make(map[string]struct{})
	for _, role := range userRoles {
		userRoleMap[role] = struct{}{}
	}

	menu := g.Menus.Menus[menuName]
	additionalData := menu.AdditionalData

	dropdownSlice := make([]discordgo.SelectMenuOption, len(additionalData.EncounterDropdown))
	copy(dropdownSlice, additionalData.EncounterDropdown)

	// set default selections based on roles present
	for i, _ := range dropdownSlice {
		option := &dropdownSlice[i]
		if _, ok := userRoleMap[option.Value]; ok {
			option.Default = true
		}
	}

	// generate role list description
	descriptionRoleList := "\n### Available roles"
	for _, role := range dropdownSlice {
		descriptionRoleList += fmt.Sprintf("\n- %s", role.Label)
	}

	minValues := 0
	maxValues := 1

	if additionalData.MultiSelect {
		maxValues = len(dropdownSlice)
	}

	// format response message
	customID := []string{string(MenuEncounter), string(CommandEncounterProcess), menuName}
	dropdownMenu := discordgo.SelectMenu{
		MenuType:  discordgo.StringSelectMenu,
		CustomID:  strings.Join(customID, " "),
		Options:   dropdownSlice,
		MinValues: &minValues,
		MaxValues: maxValues,
	}

	message := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       g.Menus.Menus[menuName].Title,
					Description: g.Menus.Menus[menuName].Description + descriptionRoleList,
				},
			},
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						dropdownMenu,
					},
				},
			},
		},
	}

	if len(g.Menus.Menus[menuName].ImageURL) != 0 {
		message.Data.Embeds[0].Image = &discordgo.MessageEmbedImage{
			URL: g.Menus.Menus[menuName].ImageURL,
		}
	}

	err := s.InteractionRespond(i.Interaction, message)
	if err != nil {
		fmt.Printf("Unable to respond with menu: %v", err)
		return
	}

}

func (c *Clearingway) MenuEncounterProcess(s *discordgo.Session, i *discordgo.InteractionCreate, menuName string) {
	err := discord.StartInteraction(s, i.Interaction, "Processing request...")
	if err != nil {
		fmt.Printf("Error sending Discord message: %v\n", err)
		return
	}

	g, ok := c.Guilds.Guilds[i.GuildID]
	if !ok {
		fmt.Printf("Interaction received from guild %s with no configuration!\n", i.GuildID)
		return
	}

	menu := g.Menus.Menus[menuName]
	additionalData := menu.AdditionalData

	userOptions := i.MessageComponentData().Values
	userOptionsMap := make(map[string]struct{})
	for _, opt := range userOptions {
		userOptionsMap[opt] = struct{}{}
	}

	userRoles := i.Member.Roles
	userRolesMap := make(map[string]struct{})
	for _, roleId := range userRoles {
		userRolesMap[roleId] = struct{}{}
	}

	// check which roles to add or remove
	rolesToAdd := []*MenuRoleHelper{}
	rolesToRemove := []*Role{}
	for _, roleHelper := range additionalData.Roles {
		role := roleHelper.Role

		_, selected := userOptionsMap[role.DiscordRole.ID]
		_, present := userRolesMap[role.DiscordRole.ID]

		if selected && !present {
			rolesToAdd = append(rolesToAdd, roleHelper)
		} else if !selected && present {
			rolesToRemove = append(rolesToRemove, role)
		}
	}

	successRoles := []string{}
	removedRoles := []string{}
	failedRoles := []string{}
	errorRoles := []string{}

	// add/remove roles based on prereq met
	for _, roleHelper := range rolesToAdd {
		prereq := roleHelper.Prerequisite
		if roleHelper.Prerequisite != nil {
			if _, prereqMet := userRolesMap[prereq.DiscordRole.ID]; !prereqMet {
				failedRoles = append(failedRoles, roleHelper.Prerequisite.DiscordRole.ID)
				continue
			}
		}

		err := roleHelper.Role.AddToCharacter(i.GuildID, i.Member.User.ID, s)
		if err != nil {
			fmt.Printf("Error adding role: %+v\n", err)
			errorRoles = append(errorRoles, roleHelper.Role.DiscordRole.ID)
			continue
		}

		successRoles = append(successRoles, roleHelper.Role.DiscordRole.ID)
	}

	for _, role := range rolesToRemove {
		err := role.RemoveFromCharacter(i.GuildID, i.Member.User.ID, s)
		if err != nil {
			fmt.Printf("Error removing role: %+v\n", err)
			errorRoles = append(errorRoles, role.DiscordRole.ID)
			continue
		}

		removedRoles = append(removedRoles, role.DiscordRole.ID)
	}

	// form response string
	responseMsg := ""

	for i, roles := range [2][]string{successRoles, removedRoles} {
		verb := ""
		if i == 0 {
			verb = "added"
		} else {
			verb = "removed"
		}

		length := len(roles)
		if length == 1 {
			responseMsg += fmt.Sprintf("Successfully %v role: <@&%v>\n", verb, roles[0])
		} else if length > 1 { // exclude case where no roles were added/removed
			responseMsg += fmt.Sprintf("Successfully %v roles: ", verb)
			for _, role := range roles {
				responseMsg += fmt.Sprintf("<@&%v> ", role)
			}
			responseMsg += "\n"
		}
	}

	for _, failedRole := range failedRoles {
		responseMsg += fmt.Sprintf("You do not have the required role: <@&%v>\n", failedRole)
	}

	if len(errorRoles) != 0 {
		responseMsg += "An error has occurred. Please contact an admin with this message.\nRole actions that failed: "
		for _, role := range errorRoles {
			responseMsg += fmt.Sprintf("<@&%v>\n", role)
		}
	}

	if len(responseMsg) == 0 {
		responseMsg = "Nothing has been done!\n"
	}
	_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &responseMsg,
	})
	if err != nil {
		fmt.Printf("Error sending Discord message: %v\n", err)
		return
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

// MenuVerifySendModal sends the user a modal that asks for their character's
// first name, last name, and world to verify their clears
func MenuVerifySendModal(s *discordgo.Session, i *discordgo.InteractionCreate) {
	customID := []string{string(MenuVerify), string(CommandClearsModal)}
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID: strings.Join(customID, " "),
			Title:    "Verify your clears",
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:    "firstName",
							Style:       discordgo.TextInputShort,
							Label:       "Character first name",
							Placeholder: "First name",
							Required:    true,
							MaxLength:   15,
							MinLength:   2,
						},
					},
				},
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:    "lastName",
							Style:       discordgo.TextInputShort,
							Label:       "Character last name",
							Placeholder: "Last name",
							Required:    true,
							MaxLength:   15,
							MinLength:   2,
						},
					},
				},
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:    "world",
							Style:       discordgo.TextInputShort,
							Label:       "Character world",
							Placeholder: "World",
							Required:    true,
							MaxLength:   20,
							MinLength:   2,
						},
					},
				},
			},
		},
	})
	if err != nil {
		fmt.Printf("Unable to send modal: %v", err)
		return
	}
}

// MenuVerifyProcess takes the submitted data from the modal above and processes the character
func (c *Clearingway) MenuVerifyProcess(s *discordgo.Session, i *discordgo.InteractionCreate) {
	g, ok := c.Guilds.Guilds[i.GuildID]
	if !ok {
		fmt.Printf("Interaction received from guild %s with no configuration!\n", i.GuildID)
		return
	}

	// Get the text fields
	options := i.ModalSubmitData()

	err := discord.StartInteraction(s, i.Interaction, "Received character info...")
	if err != nil {
		fmt.Printf("Error sending Discord message: %v\n", err)
		return
	}

	var world string
	var firstName string
	var lastName string

	firstName = options.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
	lastName = options.Components[1].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
	world = options.Components[2].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value

	c.ClearsHelper(s, i, g, world, firstName, lastName)
}
