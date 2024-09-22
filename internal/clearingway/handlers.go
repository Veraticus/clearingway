package clearingway

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/Veraticus/clearingway/internal/discord"

	"github.com/bwmarrin/discordgo"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type pendingRole struct {
	role    *Role
	message string
}

func (c *Clearingway) DiscordReady(s *discordgo.Session, event *discordgo.Ready) {
	fmt.Printf("Initializing Discord...\n")

	for _, discordGuild := range event.Guilds {
		gid := discordGuild.ID
		guild, ok := c.Guilds.Guilds[discordGuild.ID]
		if !ok {
			fmt.Printf("Initialized in guild %s with no configuration!\n", gid)
			continue
		}
		existingRoles, err := s.GuildRoles(gid)
		if err != nil {
			fmt.Printf("Error getting existing roles: %v\n", err)
			return
		}

		fmt.Printf("Initializing roles...\n")
		for _, r := range guild.AllRoles() {
			fmt.Printf("Ensuring role: %+v\n", r)
			err := r.Ensure(gid, s, existingRoles)
			if err != nil {
				fmt.Printf("Error ensuring role %+v: %+v\n", r, err)
			}
		}

		for _, menu := range guild.Menus.Menus {
			if menu.Type == MenuEncounter {
				additionalData := menu.AdditionalData
				menu.MenuEncounterInit(guild.Encounters, additionalData.RoleType)
			}
		}

		time.Sleep(1 * time.Second)

		fmt.Printf("Adding commands...")

		commandList := []*discordgo.ApplicationCommand{
			ClearCommand,
			UncomfyCommand,
			UncolorCommand,
			RemoveCommand,
			RolesCommand,
		}

		if guild.IsProgEnabled() {
			commandList = append(commandList, ProgCommand)
		}

		if guild.ReclearsEnabled {
			commandList = append(commandList, ReclearCommand)
		}

		if guild.NameColorsEnabled {
			commandList = append(commandList, NameColorCommand)
		}

		if guild.MenuEnabled {
			commandList = append(commandList, MenuCommand)
		}

		addedCommands, err := s.ApplicationCommandBulkOverwrite(event.User.ID, discordGuild.ID, commandList)
		fmt.Printf("List of successfully created commands:\n")
		for _, command := range addedCommands {
			fmt.Printf("\t%v\n", command.Name)
		}
		if err != nil {
			fmt.Printf("Could not add some commands: %v\n", err)
		}

		// fmt.Printf("Removing commands...\n")
		// cmd, err := s.ApplicationCommandCreate(event.User.ID, guild.ID, verifyCommand)
		// if err != nil {
		// 	fmt.Printf("Could not find command: %v\n", err)
		// }
		// err = s.ApplicationCommandDelete(event.User.ID, guild.ID, cmd.ID)
		// if err != nil {
		// 	fmt.Printf("Could not delete command: %v\n", err)
		// }
	}
	fmt.Printf("Clearingway ready!\n")
	c.Ready = true
}

var adminPermission int64 = discordgo.PermissionAdministrator

var MenuCommand = &discordgo.ApplicationCommand{
	Name:                     "menu",
	Description:              "Send the roles menu to the current channel.",
	DefaultMemberPermissions: &adminPermission,
}

var ClearCommand = &discordgo.ApplicationCommand{
	Name:        "clears",
	Description: "Verify you own your character and assign them cleared roles.",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:         discordgo.ApplicationCommandOptionString,
			Name:         "world",
			Description:  "Your character's world",
			Required:     true,
			Autocomplete: true,
		},
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "first-name",
			Description: "Your character's first name",
			Required:    true,
		},
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "last-name",
			Description: "Your character's last name",
			Required:    true,
		},
	},
}

var UncomfyCommand = &discordgo.ApplicationCommand{
	Name:        "uncomfy",
	Description: "Use this command to remove Comfy roles if you don't want them.",
}

var UncolorCommand = &discordgo.ApplicationCommand{
	Name:        "uncolor",
	Description: "Use this command to remove parsing roles if you don't want them.",
}

var RemoveCommand = &discordgo.ApplicationCommand{
	Name:        "removeall",
	Description: "Use this command to remove all Clearingway-related roles if you don't want them.",
}

var RolesCommand = &discordgo.ApplicationCommand{
	Name:        "roles",
	Description: "See what roles Clearingway has set up and how to get them.",
}

var ProgCommand = &discordgo.ApplicationCommand{
	Name:        "prog",
	Description: "Assign yourself roles based on prog from a linked fflogs report",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:         discordgo.ApplicationCommandOptionString,
			Name:         "world",
			Description:  "Your character's world",
			Required:     true,
			Autocomplete: true,
		},
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "first-name",
			Description: "Your character's first name",
			Required:    true,
		},
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "last-name",
			Description: "Your character's last name",
			Required:    true,
		},
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "report-id",
			Description: "An fflogs report URL or ID",
			Required:    true,
		},
	},
}

var ReclearCommand = &discordgo.ApplicationCommand{
	Name:        "reclears",
	Description: "Assign or remove reclear roles",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "ultimate",
			Description: "The ultimate you want to reclear",
			Required:    true,
			Choices: []*discordgo.ApplicationCommandOptionChoice{
				{
					Name:  "UCoB",
					Value: "The Unending Coil of Bahamut (Ultimate)",
				},
				{
					Name:  "UWU",
					Value: "The Weapon's Refrain (Ultimate)",
				},
				{
					Name:  "TEA",
					Value: "The Epic of Alexander (Ultimate)",
				},
				{
					Name:  "DSR",
					Value: "Dragonsong's Reprise (Ultimate)",
				},
				{
					Name:  "TOP",
					Value: "The Omega Protocol (Ultimate)",
				},
				// TODO: implement when FRU goes live
				/*
					{
						Name: "FRU",
						Value: "Futures Rewritten (Ultimate)",
					},
				*/
			},
		},
	},
}

var NameColorCommand = &discordgo.ApplicationCommand{
	Name:        "namecolor",
	Description: "Assign or remove name color roles",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "ultimate",
			Description: "The ultimate that corresponds to the color you want",
			Required:    true,
			Choices: []*discordgo.ApplicationCommandOptionChoice{
				{
					Name:  "UCoB",
					Value: "The Unending Coil of Bahamut (Ultimate)",
				},
				{
					Name:  "UWU",
					Value: "The Weapon's Refrain (Ultimate)",
				},
				{
					Name:  "TEA",
					Value: "The Epic of Alexander (Ultimate)",
				},
				{
					Name:  "DSR",
					Value: "Dragonsong's Reprise (Ultimate)",
				},
				{
					Name:  "TOP",
					Value: "The Omega Protocol (Ultimate)",
				},
				// TODO: implement when FRU goes live
				/*
					{
						Name: "FRU",
						Value: "Futures Rewritten (Ultimate)",
					},
				*/
			},
		},
	},
}

func (c *Clearingway) InteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		switch i.ApplicationCommandData().Name {
		case "clears":
			c.Clears(s, i)
		case "uncomfy":
			c.Uncomfy(s, i)
		case "uncolor":
			c.Uncolor(s, i)
		case "roles":
			c.Roles(s, i)
		case "prog":
			c.Prog(s, i)
		case "removeall":
			c.RemoveAll(s, i)
		case "namecolor":
			c.ToggleColor(s, i)
		case "reclears":
			c.ToggleReclear(s, i)
		case "menu":
			c.MenuMainSend(s, i)
		}
	case discordgo.InteractionApplicationCommandAutocomplete:
		c.Autocomplete(s, i)
	case discordgo.InteractionMessageComponent:
		customID := i.MessageComponentData().CustomID
		command := strings.Split(customID, " ")
		if ok := len(command) > 1; !ok {
			fmt.Printf("Invalid custom ID received: \"%v\"\n", customID)
			return
		}

		switch MenuType(command[0]) {
		case MenuVerify:
			switch CommandType(command[1]) {
			case CommandMenu:
				// send_modal()
			case CommandClearsModal:
				// Clears()
			}
		case MenuRemove:
			switch CommandType(command[1]) {
			case CommandMenu:
				// send_menu()
			case CommandRemoveComfy:
				// Uncomfy()
			case CommandRemoveColor:
				// Uncolor()
			case CommandRemoveAll:
				// RemoveAll()
			}
		case MenuEncounter:
			if ok := len(command) > 2; !ok {
				fmt.Printf("Invalid custom ID received: \"%v\"\n", customID)
				return
			}
			switch CommandType(command[1]) {
			case CommandMenu:
				c.MenuEncounterSend(s, i, command[2])
			case CommandEncounterProcess:
				c.MenuEncounterProcess(s, i, command[2])
			}
		}
	}
}

func (c *Clearingway) Uncomfy(s *discordgo.Session, i *discordgo.InteractionCreate) {
	g, ok := c.Guilds.Guilds[i.GuildID]
	if !ok {
		fmt.Printf("Interaction received from guild %s with no configuration!\n", i.GuildID)
		return
	}

	// Ignore messages not on the correct channel
	if i.ChannelID != g.ChannelId {
		fmt.Printf("Ignoring message not in channel %s.\n", g.ChannelId)
	}

	err := discord.StartInteraction(s, i.Interaction, "Uncomfying you...")
	if err != nil {
		fmt.Printf("Error sending Discord message: %v\n", err)
		return
	}

	member, err := c.Discord.Session.GuildMember(g.Id, i.Member.User.ID)
	if err != nil {
		err = discord.ContinueInteraction(s, i.Interaction, err.Error())
		if err != nil {
			fmt.Printf("Error sending Discord message: %v\n", err)
		}
		return
	}

	uncomfyRoles := []*Role{}
	for _, r := range g.AllRoles() {
		if r.Skip {
			continue
		}
		if r.Uncomfy {
			uncomfyRoles = append(uncomfyRoles, r)
		}
	}

	if len(uncomfyRoles) == 0 {
		err = discord.ContinueInteraction(s, i.Interaction, "Uncomfy roles are not present in this Discord!")
		if err != nil {
			fmt.Printf("Error sending Discord message: %v\n", err)
		}
		return
	}

	rolesToRemove := []*Role{}
	for _, r := range uncomfyRoles {
		if r.DiscordRole == nil {
			fmt.Printf("Cannot uncomfy %+v as it has not connected to a Discord role!\n", r)
			continue
		}
		if r.PresentInRoles(member.Roles) {
			rolesToRemove = append(rolesToRemove, r)
		}
	}
	if len(rolesToRemove) == 0 {
		err = discord.ContinueInteraction(s, i.Interaction, "You do not have any uncomfy roles!")
		if err != nil {
			fmt.Printf("Error sending Discord message: %v\n", err)
		}
		return
	}

	for _, r := range rolesToRemove {
		err = r.RemoveFromCharacter(g.Id, i.Member.User.ID, c.Discord.Session)
		if err != nil {
			fmt.Printf("Error removing uncomfy role: %+v\n", err)
		}
		fmt.Printf("Removing uncomfy role: %+v\n", err)
	}

	err = discord.ContinueInteraction(s, i.Interaction, "_ _\n__Uncomfy roles:__\n⮕ Removed!\n")
	if err != nil {
		fmt.Printf("Error sending Discord message: %v\n", err)
	}
}

func (c *Clearingway) Uncolor(s *discordgo.Session, i *discordgo.InteractionCreate) {
	g, ok := c.Guilds.Guilds[i.GuildID]
	if !ok {
		fmt.Printf("Interaction received from guild %s with no configuration!\n", i.GuildID)
		return
	}

	// Ignore messages not on the correct channel
	if i.ChannelID != g.ChannelId {
		fmt.Printf("Ignoring message not in channel %s.\n", g.ChannelId)
	}

	err := discord.StartInteraction(s, i.Interaction, "Uncoloring you...")
	if err != nil {
		fmt.Printf("Error sending Discord message: %v\n", err)
		return
	}

	member, err := c.Discord.Session.GuildMember(g.Id, i.Member.User.ID)
	if err != nil {
		err = discord.ContinueInteraction(s, i.Interaction, err.Error())
		if err != nil {
			fmt.Printf("Error sending Discord message: %v\n", err)
		}
		return
	}

	uncolorRoles := []*Role{}
	for _, r := range g.AllRoles() {
		if r.Skip {
			continue
		}
		if r.Uncolor {
			uncolorRoles = append(uncolorRoles, r)
		}
	}

	if len(uncolorRoles) == 0 {
		err = discord.ContinueInteraction(s, i.Interaction, "Parsing roles are not present in this Discord!")
		if err != nil {
			fmt.Printf("Error sending Discord message: %v\n", err)
		}
		return
	}

	rolesToRemove := []*Role{}
	for _, r := range uncolorRoles {
		if r.PresentInRoles(member.Roles) {
			rolesToRemove = append(rolesToRemove, r)
		}
	}
	if len(rolesToRemove) == 0 {
		err = discord.ContinueInteraction(s, i.Interaction, "You do not have any parsing roles!")
		if err != nil {
			fmt.Printf("Error sending Discord message: %v\n", err)
		}
		return
	}

	for _, r := range rolesToRemove {
		err = r.RemoveFromCharacter(g.Id, i.Member.User.ID, c.Discord.Session)
		if err != nil {
			fmt.Printf("Error removing parsing role: %+v\n", err)
		}
		fmt.Printf("Removing parsing role: %+v\n", err)
	}

	err = discord.ContinueInteraction(s, i.Interaction, "_ _\n__Parsing roles:__\n⮕ Removed!\n")
	if err != nil {
		fmt.Printf("Error sending Discord message: %v\n", err)
	}
}

func (c *Clearingway) Roles(s *discordgo.Session, i *discordgo.InteractionCreate) {
	g, ok := c.Guilds.Guilds[i.GuildID]
	if !ok {
		fmt.Printf("Interaction received from guild %s with no configuration!\n", i.GuildID)
		return
	}

	// Ignore messages not on the correct channel
	if i.ChannelID != g.ChannelId {
		fmt.Printf("Ignoring message not in channel %s.\n", g.ChannelId)
	}

	chunks := discord.NewChunks()
	chunks.Write("_ _\n")
	chunks.Write("Clearingway considers the following encounters relevant:\n")

	for _, e := range g.Encounters.Encounters {
		chunks.Write("⮕ " + e.Name + "\n")
	}

	chunks.Write("\nClearingway can give out these roles with `/clears`:\n")

	for _, r := range g.AllRoles() {
		if r.Skip {
			continue
		}
		chunks.Write(fmt.Sprintf("__**%s**__\n⮕ %s\n\n", r.Name, r.Description))
	}

	for n, c := range chunks.Chunks {
		var err error
		if n == 0 {
			err = discord.StartInteraction(s, i.Interaction, c.String())
		} else {
			err = discord.ContinueInteraction(s, i.Interaction, c.String())
		}
		if err != nil {
			fmt.Printf("Error sending Discord message: %v\n", err)
			return
		}
	}
}

func (c *Clearingway) RemoveAll(s *discordgo.Session, i *discordgo.InteractionCreate) {
	g, ok := c.Guilds.Guilds[i.GuildID]
	if !ok {
		fmt.Printf("Interaction received from guild %s with no configuration!\n", i.GuildID)
		return
	}

	// Ignore messages not on the correct channel
	if i.ChannelID != g.ChannelId {
		fmt.Printf("Ignoring message not in channel %s.\n", g.ChannelId)
	}

	err := discord.StartInteraction(s, i.Interaction, "Removing all Clearingway related roles...")
	if err != nil {
		fmt.Printf("Error sending Discord message: %v\n", err)
		return
	}

	member, err := c.Discord.Session.GuildMember(g.Id, i.Member.User.ID)
	if err != nil {
		err = discord.ContinueInteraction(s, i.Interaction, err.Error())
		if err != nil {
			fmt.Printf("Error sending Discord message: %v\n", err)
		}
		return
	}

	// Get a list of all Clearingway-related roles configured for the current guild, excludes roles with the Skip flag
	// List is used to check which roles to remove from the user
	clearingwayRoles := []*Role{}
	for _, r := range g.AllRoles() {
		if r.Skip {
			continue
		}
		clearingwayRoles = append(clearingwayRoles, r)
	}

	if len(clearingwayRoles) == 0 {
		err = discord.ContinueInteraction(s, i.Interaction, "No Clearingway related roles are present in this Discord!")
		if err != nil {
			fmt.Printf("Error sending Discord message: %v\n", err)
		}
		return
	}

	rolesToRemove := []*Role{}
	for _, r := range clearingwayRoles {
		if r.DiscordRole == nil {
			fmt.Printf("Cannot remove %+v as it has not connected to a Discord role!\n", r)
			continue
		}
		if r.PresentInRoles(member.Roles) {
			rolesToRemove = append(rolesToRemove, r)
		}
	}
	if len(rolesToRemove) == 0 {
		err = discord.ContinueInteraction(s, i.Interaction, "You do not have any Clearingway-related roles!")
		if err != nil {
			fmt.Printf("Error sending Discord message: %v\n", err)
		}
		return
	}

	for _, r := range rolesToRemove {
		err = r.RemoveFromCharacter(g.Id, i.Member.User.ID, c.Discord.Session)
		if err != nil {
			fmt.Printf("Error removing role: %+v\n", err)
		}
		fmt.Printf("Removing role: %+v\n", r.Name)
	}

	err = discord.ContinueInteraction(s, i.Interaction, "_ _\n__Clearingway-related roles:__\n⮕ Removed!\n")
	if err != nil {
		fmt.Printf("Error sending Discord message: %v\n", err)
	}
}

func (c *Clearingway) ToggleReclear(s *discordgo.Session, i *discordgo.InteractionCreate) {
	g, ok := c.Guilds.Guilds[i.GuildID]
	if !ok {
		fmt.Printf("Interaction received from guild %s with no configuration!\n", i.GuildID)
		return
	}

	err := discord.StartInteraction(s, i.Interaction, "Checking for respective clear role...")
	if err != nil {
		fmt.Printf("Error sending Discord message: %v\n", err)
		return
	}

	ultimate := i.ApplicationCommandData().Options[0].StringValue()
	encounter := g.Encounters.ForName(ultimate)
	reclearRole := encounter.Roles[ReclearRole]
	clearedRole := encounter.Roles[ClearedRole]

	var rolePresent = false
	var clearPresent = false
	for _, r := range i.Member.Roles {
		if r == reclearRole.DiscordRole.ID {
			rolePresent = true
			continue
		}
		if r == clearedRole.DiscordRole.ID {
			clearPresent = true
			continue
		}
	}

	// Remove role no matter what
	// Add role only if cleared role is present
	if rolePresent {
		fmt.Printf("Removing role: %+v\n", reclearRole.Name)
		err = reclearRole.RemoveFromCharacter(g.Id, i.Member.User.ID, c.Discord.Session)
		if err != nil {
			fmt.Printf("Error removing role: %+v\n", err)
			return
		}
		tempstr := fmt.Sprintf("Successfully removed role: <@&%v>", reclearRole.DiscordRole.ID)
		err = discord.ContinueInteraction(s, i.Interaction, tempstr)
		if err != nil {
			fmt.Printf("Error sending Discord message: %v\n", err)
		}
	} else {
		if clearPresent {
			fmt.Printf("Adding role: %+v\n", reclearRole.Name)
			err = reclearRole.AddToCharacter(g.Id, i.Member.User.ID, c.Discord.Session)
			if err != nil {
				fmt.Printf("Error adding role: %+v\n", err)
				return
			}
			tempstr := fmt.Sprintf("Successfully added role: <@&%v>", reclearRole.DiscordRole.ID)
			err = discord.ContinueInteraction(s, i.Interaction, tempstr)
			if err != nil {
				fmt.Printf("Error sending Discord message: %v\n", err)
			}
		} else {
			tempstr := fmt.Sprintf("You do not have the required role: <@&%v>", clearedRole.DiscordRole.ID)
			err = discord.ContinueInteraction(s, i.Interaction, tempstr)
			if err != nil {
				fmt.Printf("Error sending Discord message: %v\n", err)
			}
		}
	}
}

func (c *Clearingway) ToggleColor(s *discordgo.Session, i *discordgo.InteractionCreate) {
	g, ok := c.Guilds.Guilds[i.GuildID]
	if !ok {
		fmt.Printf("Interaction received from guild %s with no configuration!\n", i.GuildID)
		return
	}

	err := discord.StartInteraction(s, i.Interaction, "Checking for respective clear role...")
	if err != nil {
		fmt.Printf("Error sending Discord message: %v\n", err)
		return
	}

	wantedUltimate := i.ApplicationCommandData().Options[0].StringValue()

	// Get all possible roles and the encounter while we're at it
	var colorRoles []*Role
	var wantedEncounter *Encounter
	for _, e := range UltimateEncounters.Encounters {
		encounter := g.Encounters.ForName(e.Name)
		colorRoles = append(colorRoles, encounter.Roles[ColorRole])
		if e.Name == wantedUltimate {
			wantedEncounter = encounter
		}
	}

	requestedColorRole := wantedEncounter.Roles[ColorRole]
	clearedRole := wantedEncounter.Roles[ClearedRole]

	var roleToRemove *Role
	var cleared = false
	for _, memberRole := range i.Member.Roles {
		if roleToRemove == nil {
			for _, colorRole := range colorRoles {
				if colorRole.DiscordRole.ID == memberRole {
					roleToRemove = colorRole
					break
				}
			}
		}
		if (!cleared) && (memberRole == clearedRole.DiscordRole.ID) {
			cleared = true
		}
	}

	if cleared {
		// remove existing color role
		tempstr := ""
		if roleToRemove != nil {
			err = removeRoleHelper(s, i.Interaction, roleToRemove)
			if err != nil {
				return
			}
			tempstr += fmt.Sprintf("Successfully removed role: <@&%v>", roleToRemove.DiscordRole.ID)
		}
		// add role if requested role is not the same as color role
		if requestedColorRole != roleToRemove {
			err = addRoleHelper(s, i.Interaction, requestedColorRole)
			if err != nil {
				return
			}
			tempstr += fmt.Sprintf("\nSuccessfully added role: <@&%v>", requestedColorRole.DiscordRole.ID)
		}
		discord.ContinueInteraction(s, i.Interaction, tempstr)
	} else {
		// remove the requested color role if it's the same
		// edge case where someone has clears removed and doesn't want the color
		if requestedColorRole == roleToRemove {
			err = removeRoleHelper(s, i.Interaction, roleToRemove)
			if err != nil {
				return
			}
			tempstr := fmt.Sprintf("Successfully removed role: <@&%v>", roleToRemove.DiscordRole.ID)
			discord.ContinueInteraction(s, i.Interaction, tempstr)
		} else {
			// user doesn't meet the requirements
			tempstr := fmt.Sprintf("You do not have the required role: <@&%v>", clearedRole.DiscordRole.ID)
			err = discord.ContinueInteraction(s, i.Interaction, tempstr)
			if err != nil {
				fmt.Printf("Error sending Discord message: %v\n", err)
			}
		}
	}

}

func removeRoleHelper(s *discordgo.Session, i *discordgo.Interaction, roleToRemove *Role) error {
	fmt.Printf("Removing role: %+v\n", roleToRemove.Name)
	err := roleToRemove.RemoveFromCharacter(i.GuildID, i.Member.User.ID, s)
	if err != nil {
		fmt.Printf("Error removing role: %+v\n", err)
		return err
	}
	return nil
}

func addRoleHelper(s *discordgo.Session, i *discordgo.Interaction, roleToAdd *Role) error {
	fmt.Printf("Adding role: %+v\n", roleToAdd.Name)
	err := roleToAdd.AddToCharacter(i.GuildID, i.Member.User.ID, s)
	if err != nil {
		fmt.Printf("Error adding role: %+v\n", err)
		return err
	}
	return nil
}

func (c *Clearingway) Autocomplete(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	var world string
	if option, ok := optionMap["world"]; ok {
		world = option.StringValue()
	}

	choices := []*discordgo.ApplicationCommandOptionChoice{}
	title := cases.Title(language.AmericanEnglish)

	if len(world) == 0 {
		for _, world := range c.AllWorlds {
			worldTitle := title.String(world)
			choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
				Name:  worldTitle,
				Value: world,
			})
		}
		return
	} else {
		for _, worldCompletion := range c.AutoCompleteTrie.SearchAll(world) {
			worldCompletionTitle := title.String(worldCompletion)
			choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
				Name:  worldCompletionTitle,
				Value: worldCompletion,
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

var nonWorldRegex = regexp.MustCompile(`[^a-zA-Z]+`)

func cleanWorld(world string) string {
	return nonWorldRegex.ReplaceAllString(world, "")
}
