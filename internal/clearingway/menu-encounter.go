package clearingway

import (
	"fmt"
	"strings"

	"github.com/Veraticus/clearingway/internal/discord"
	"github.com/bwmarrin/discordgo"
)

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
		descriptionRoleList += fmt.Sprintf("\n- <@&%s>", role.Value)
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

	menu = g.Menus.Menus[menuName]

	message := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       menu.Title,
					Description: menu.Description + descriptionRoleList,
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

	menu.MenuStyle(message.Data.Embeds)

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
