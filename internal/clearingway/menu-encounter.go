package clearingway

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/Veraticus/clearingway/internal/discord"
	"github.com/bwmarrin/discordgo"
)

const DEFAULT_TITLE = "General"

func (m *Menu) MenuEncounterInit(es *Encounters, roleTypes []RoleType) {
	additionalData := m.AdditionalData
	additionalData.Roles = make(map[string]*MenuRoleHelper)
	dropdownSlice := []DropdownHelper{}
	difficultyToIndex := make(map[string]int)

	// populate individual SelectMenus
	if len(additionalData.Difficulties) > 0 {
		for index, difficulty := range additionalData.Difficulties {
			dropdownSlice = append(dropdownSlice, DropdownHelper{
				DropdownTitle:     difficulty,
				SelectMenuOptions: []discordgo.SelectMenuOption{},
			})
			difficultyToIndex[difficulty] = index
		}
	} else if len(roleTypes) > 0 {
		// if no difficulties specified but role types are specified, populate a single SelectMenu
		dropdownSlice = append(dropdownSlice, DropdownHelper{
			DropdownTitle:     DEFAULT_TITLE,
			SelectMenuOptions: []discordgo.SelectMenuOption{},
		})
		difficultyToIndex[DEFAULT_TITLE] = 0
	}

	// generate options based on encounters
	for _, roleType := range roleTypes {
		for _, encounter := range es.Encounters {
			// skip encounter if menu does not include encounter's difficulty when explicitly specified
			if len(additionalData.Difficulties) > 0 && !slices.Contains(additionalData.Difficulties, encounter.Difficulty) {
				continue
			}

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

			// if difficulties are specified, put them into their respective selectMenu
			if len(additionalData.Difficulties) != 0 {
				difficultyIndex := difficultyToIndex[encounter.Difficulty]
				menuRoleHelper.DifficultyIndex = difficultyIndex
				dropdownToAppend := &dropdownSlice[difficultyIndex]
				dropdownToAppend.SelectMenuOptions = append(dropdownToAppend.SelectMenuOptions, dropdownOption)
			} else {
				dropdownSlice[0].SelectMenuOptions = append(dropdownSlice[0].SelectMenuOptions, dropdownOption)
			}

			additionalData.Roles[role.DiscordRole.ID] = menuRoleHelper
		}
	}

	// if there are additional roles, make a new SelectMenu for it
	extraRolesIndex := 0
	if len(additionalData.ExtraRoles) != 0 {
		dropdownSlice = append(dropdownSlice, DropdownHelper{
			DropdownTitle:     "Others",
			SelectMenuOptions: []discordgo.SelectMenuOption{},
		})
		extraRolesIndex = len(dropdownSlice) - 1
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

		dropdownSlice[extraRolesIndex].SelectMenuOptions = append(dropdownSlice[extraRolesIndex].SelectMenuOptions, dropdownOption)
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

	menu, ok := g.Menus.Menus[menuName]
	if !ok {
		discord.StartInteraction(s, i.Interaction, "Error: Menu not found.")
		return
	}
	additionalData := menu.AdditionalData

	dropdowns := make([]DropdownHelper, len(additionalData.EncounterDropdown))
	copy(dropdowns, additionalData.EncounterDropdown)

	// set default selections based on roles present
	for index, dropdown := range dropdowns {
		for selectMenuOption, _ := range dropdown.SelectMenuOptions {
			option := &dropdowns[index].SelectMenuOptions[selectMenuOption]
			if _, ok := userRoleMap[option.Value]; ok {
				option.Default = true
			}
		}
	}

	// generate role list descriptions
	descriptionRoleList := []string{}
	for index, dropdown := range dropdowns {
		descriptionRoleList = append(descriptionRoleList, "")
		for _, role := range dropdown.SelectMenuOptions {
			descriptionRoleList[index] += fmt.Sprintf("\n- <@&%s>", role.Value)
		}
	}

	minValues := 0
	maxValues := []int{}

	if additionalData.MultiSelect {
		for _, dropdown := range dropdowns {
			maxValues = append(maxValues, len(dropdown.SelectMenuOptions))
		}
	} else {
		for _, _ = range dropdowns {
			maxValues = append(maxValues, 1)
		}
	}

	// format response message
	baseCustomID := []string{string(MenuEncounter), string(CommandEncounterProcess), menuName}
	fields := []*discordgo.MessageEmbedField{}
	components := []discordgo.MessageComponent{}

	for index, dropdown := range dropdowns {
		// add individual SelectMenus to message
		dropdownCustomID := append(baseCustomID, strconv.Itoa(index))
		components = append(components, discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.SelectMenu{
					MenuType:  discordgo.StringSelectMenu,
					CustomID:  strings.Join(dropdownCustomID, " "),
					Options:   dropdown.SelectMenuOptions,
					MinValues: &minValues,
					MaxValues: maxValues[index],
				},
			},
		})

		// if there's more than one dropdown, add the role list descriptions to a field
		if len(dropdowns) > 1 {
			fields = append(fields, &discordgo.MessageEmbedField{
				Name:   dropdown.DropdownTitle,
				Value:  descriptionRoleList[index],
				Inline: true,
			})
		}
	}

	menu = g.Menus.Menus[menuName]
	description := menu.Description + "\n### Available roles"

	embed := []*discordgo.MessageEmbed{
		{
			Title:       menu.Title,
			Description: description,
		},
	}

	// if there's only one dropdown, add the role list descriptions to the main description
	if len(dropdowns) == 1 {
		embed[0].Description += descriptionRoleList[0]
	} else {
		embed[0].Fields = fields
	}

	message := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:      discordgo.MessageFlagsEphemeral,
			Embeds:     embed,
			Components: components,
		},
	}

	menu.MenuStyle(message.Data.Embeds)

	err := s.InteractionRespond(i.Interaction, message)
	if err != nil {
		fmt.Printf("Unable to respond with menu: %v", err)
		return
	}

}

func (c *Clearingway) MenuEncounterProcess(s *discordgo.Session, i *discordgo.InteractionCreate, menuName string, menuIndexString string) {
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

	menu, ok := g.Menus.Menus[menuName]
	if !ok {
		discord.StartInteraction(s, i.Interaction, "Error: Menu not found.")
		return
	}
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

	// menuIndex specifies the index of the SelectMenu
	// indicates what roles are relevant
	menuIndex, err := strconv.Atoi(menuIndexString)
	if err != nil {
		discord.StartInteraction(s, i.Interaction, "Error: Invalid custom ID.")
		return
	}

	rolesToAdd := []*MenuRoleHelper{}
	rolesToRemove := []*Role{}
	if menuIndex < 0 {
		// menuIndex of -1 = remove roles buttons
		for _, roleHelper := range additionalData.Roles {
			role := roleHelper.Role

			_, present := userRolesMap[role.DiscordRole.ID]

			if present {
				rolesToRemove = append(rolesToRemove, role)
			}
		}
	} else {
		// check which roles to add or remove
		rolesInDifficulty := additionalData.EncounterDropdown[menuIndex]
		for _, selectMenuOption := range rolesInDifficulty.SelectMenuOptions {
			roleID := selectMenuOption.Value
			roleHelper := additionalData.Roles[roleID]
			role := roleHelper.Role

			_, selected := userOptionsMap[role.DiscordRole.ID]
			_, present := userRolesMap[role.DiscordRole.ID]

			if selected && !present {
				rolesToAdd = append(rolesToAdd, roleHelper)
			} else if !selected && present {
				rolesToRemove = append(rolesToRemove, role)
			}
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
