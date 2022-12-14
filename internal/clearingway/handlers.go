package clearingway

import (
	"fmt"
	"time"

	"github.com/Veraticus/clearingway/internal/discord"
	"github.com/Veraticus/clearingway/internal/fflogs"
	"github.com/Veraticus/clearingway/internal/ffxiv"
	"github.com/Veraticus/clearingway/internal/lodestone"

	"github.com/bwmarrin/discordgo"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

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

		fmt.Printf("Adding commands...\n")
		_, err = s.ApplicationCommandCreate(event.User.ID, discordGuild.ID, ClearCommand)
		if err != nil {
			fmt.Printf("Could not add clears command: %v\n", err)
		}

		_, err = s.ApplicationCommandCreate(event.User.ID, discordGuild.ID, UncomfyCommand)
		if err != nil {
			fmt.Printf("Could not add uncomfy command: %v\n", err)
		}

		_, err = s.ApplicationCommandCreate(event.User.ID, discordGuild.ID, RolesCommand)
		if err != nil {
			fmt.Printf("Could not add uncomfy command: %v\n", err)
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

var RolesCommand = &discordgo.ApplicationCommand{
	Name:        "roles",
	Description: "See what roles Clearingway has set up and how to get them.",
}

func (c *Clearingway) InteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		switch i.ApplicationCommandData().Name {
		case "clears":
			c.Clears(s, i)
		case "uncomfy":
			c.Uncomfy(s, i)
		case "roles":
			c.Roles(s, i)
		}
	case discordgo.InteractionApplicationCommandAutocomplete:
		c.Autocomplete(s, i)
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
		r.RemoveFromCharacter(g.Id, i.Member.User.ID, c.Discord.Session)
		fmt.Printf("Error removing uncomfy role: %+v\n", err)
	}

	err = discord.ContinueInteraction(s, i.Interaction, "_ _\n__Uncomfy roles:__\n⮕ Removed!\n")
	if err != nil {
		fmt.Printf("Error sending Discord message: %v\n", err)
	}
	return
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
		if r.Skip == true {
			continue
		}
		chunks.Write(fmt.Sprintf("__**%s**__\n⮕ %s\n\n", r.Name, r.Description))
	}

	for n, c := range chunks.Chunks {
		var err error
		if n == 1 {
			err = discord.StartInteraction(s, i.Interaction, c.String())
		} else {
			err = discord.ContinueInteraction(s, i.Interaction, c.String())
		}
		if err != nil {
			fmt.Printf("Error sending Discord message: %v\n", err)
			return
		}
	}
	return
}

func (c *Clearingway) Clears(s *discordgo.Session, i *discordgo.InteractionCreate) {
	g, ok := c.Guilds.Guilds[i.GuildID]
	if !ok {
		fmt.Printf("Interaction received from guild %s with no configuration!\n", i.GuildID)
		return
	}

	// Ignore messages not on the correct channel
	if i.ChannelID != g.ChannelId {
		fmt.Printf("Ignoring message not in channel %s.\n", g.ChannelId)
	}

	// Retrieve all the options sent to the command
	options := i.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	var world string
	var firstName string
	var lastName string

	if option, ok := optionMap["world"]; ok {
		world = option.StringValue()
	}
	if option, ok := optionMap["first-name"]; ok {
		firstName = option.StringValue()
	}
	if option, ok := optionMap["last-name"]; ok {
		lastName = option.StringValue()
	}

	if len(world) == 0 || len(firstName) == 0 || len(lastName) == 0 {
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "`/clears` command failed! Make sure you input your world, first name, and last name.",
			},
		})
		if err != nil {
			fmt.Printf("Error sending Discord message: %v\n", err)
		}
		return
	}

	title := cases.Title(language.AmericanEnglish)
	world = title.String(world)
	firstName = title.String(firstName)
	lastName = title.String(lastName)

	err := discord.StartInteraction(s, i.Interaction,
		fmt.Sprintf("Finding `%s %s (%s)` in the Lodestone...", firstName, lastName, world),
	)
	if err != nil {
		fmt.Printf("Error sending Discord message: %v\n", err)
		return
	}

	char, err := g.Characters.Init(world, firstName, lastName)
	if err != nil {
		err = discord.ContinueInteraction(s, i.Interaction, err.Error())
		if err != nil {
			fmt.Printf("Error sending Discord message: %v\n", err)
		}
		return
	}

	err = c.Fflogs.SetCharacterLodestoneID(char)
	if err != nil {
		err = discord.ContinueInteraction(s, i.Interaction,
			fmt.Sprintf(
				"Error finding this character's Lodestone ID from FF Logs: %v\nTo make lookups faster in the future, please link your character in FF Logs to the Lodestone here: https://www.fflogs.com/lodestone/import",
				err,
			))
		if err != nil {
			fmt.Printf("Error sending Discord message: %v\n", err)
			return
		}
		err = lodestone.SetCharacterLodestoneID(char)
		if err != nil {
			err = discord.ContinueInteraction(s, i.Interaction,
				fmt.Sprintf(
					"Error finding this character's Lodestone ID in the Lodestone: %v\nIf your character name is short or very common this can frequently fail. Please link your character in FF Logs to the Lodestone here: https://www.fflogs.com/lodestone/import",
					err,
				))
			if err != nil {
				fmt.Printf("Error sending Discord message: %v\n", err)
				return
			}
		}
	}

	err = discord.ContinueInteraction(s, i.Interaction,
		fmt.Sprintf("Verifying ownership of `%s (%s)`...", char.Name(), char.World),
	)
	if err != nil {
		fmt.Printf("Error sending Discord message: %v\n", err)
	}

	discordId := i.Member.User.ID
	isOwner, err := lodestone.CharacterIsOwnedByDiscordUser(char, discordId)
	if err != nil {
		err = discord.ContinueInteraction(s, i.Interaction, err.Error())
		if err != nil {
			fmt.Printf("Error sending Discord message: %v\n", err)
		}
		return
	}
	if !isOwner {
		discord.ContinueInteraction(s, i.Interaction,
			fmt.Sprintf(
				"I could not verify your ownership of `%s (%s)`!\nIf this is your character, add the following code to your Lodestone profile and then run `/clears` again:\n\n**%s**\n\nYou can edit your Lodestone profile at https://na.finalfantasyxiv.com/lodestone/my/setting/profile/",
				char.Name(),
				char.World,
				char.LodestoneSlug(discordId),
			),
		)
		if err != nil {
			fmt.Printf("Error sending Discord message: %v\n", err)
		}
		return
	}

	err = discord.ContinueInteraction(s, i.Interaction,
		fmt.Sprintf("Analyzing logs for `%s (%s)`...", char.Name(), char.World),
	)
	if err != nil {
		fmt.Printf("Error sending Discord message: %v\n", err)
	}

	if char.UpdatedRecently() {
		err = discord.ContinueInteraction(s, i.Interaction,
			fmt.Sprintf("Finished analysis for `%s (%s)`.", char.Name(), char.World),
		)
		if err != nil {
			fmt.Printf("Error sending Discord message: %v\n", err)
		}
		return
	}

	roleTexts, err := c.UpdateCharacterInGuild(char, i.Member.User.ID, g)
	if err != nil {
		err = discord.ContinueInteraction(s, i.Interaction,
			fmt.Sprintf("Could not analyze clears for `%s (%s)`: %s", char.Name(), char.World, err),
		)
		return
	}

	chunks := discord.NewChunks()
	chunks.Write(fmt.Sprintf("Finished analysis for `%s (%s)`.\n\n", char.Name(), char.World))

	for _, roleText := range roleTexts {
		chunks.Write(roleText + "\n")
	}

	for _, c := range chunks.Chunks {
		err = discord.ContinueInteraction(s, i.Interaction, "_ _\n"+c.String())
		if err != nil {
			fmt.Printf("Error sending Discord message: %v\n", err)
		}
	}

	return
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

type pendingRole struct {
	role    *Role
	message string
}

func (c *Clearingway) UpdateCharacterInGuild(char *ffxiv.Character, discordUserId string, guild *Guild) ([]string, error) {
	rankingsToGet := []*fflogs.RankingToGet{}
	for _, encounter := range guild.AllEncounters() {
		rankingsToGet = append(rankingsToGet, &fflogs.RankingToGet{IDs: encounter.Ids, Difficulty: encounter.DifficultyInt()})
	}
	rankings, err := c.Fflogs.GetRankingsForCharacter(rankingsToGet, char)
	if err != nil {
		return nil, fmt.Errorf("Error retrieving encounter rankings: %w", err)
	}

	fmt.Printf("Found the following relevant rankings for %s (%s)...\n", char.Name(), char.World)
	for _, e := range guild.Encounters.Encounters {
		fmt.Printf("%s:\n", e.Name)
		for _, r := range e.Ranks(rankings) {
			fmt.Printf("  %+v\n", r)
		}
	}

	fmt.Printf("Found the following ultimate rankings for %s (%s)...\n", char.Name(), char.World)
	for _, e := range UltimateEncounters.Encounters {
		fmt.Printf("%s:\n", e.Name)
		for _, r := range e.Ranks(rankings) {
			fmt.Printf("  %+v\n", r)
		}
	}

	member, err := c.Discord.Session.GuildMember(guild.Id, discordUserId)
	if err != nil {
		return nil, fmt.Errorf("Could not retrieve roles for user: %w", err)
	}

	text := []string{}

	shouldApplyOpts := &ShouldApplyOpts{
		Character: char,
		Rankings:  rankings,
	}

	rolesToApply := []*pendingRole{}
	rolesToRemove := []*pendingRole{}

	// Do not include ultimate encounters for encounter, parsing,
	// and world roles, since we don't want clears for those fight
	// to count towards clears or colors.
	for _, role := range guild.NonUltRoles() {
		if role.ShouldApply == nil {
			continue
		}

		shouldApplyOpts.Encounters = guild.Encounters

		shouldApply, message := role.ShouldApply(shouldApplyOpts)
		if shouldApply {
			rolesToApply = append(rolesToApply, &pendingRole{role: role, message: message})
		} else {
			rolesToRemove = append(rolesToRemove, &pendingRole{role: role, message: message})
		}
	}

	// Add ultimate roles too
	for _, role := range guild.UltRoles() {
		if role.ShouldApply == nil {
			continue
		}

		shouldApplyOpts.Encounters = UltimateEncounters

		shouldApply, message := role.ShouldApply(shouldApplyOpts)
		if shouldApply {
			rolesToApply = append(rolesToApply, &pendingRole{role: role, message: message})
		} else {
			rolesToRemove = append(rolesToRemove, &pendingRole{role: role, message: message})
		}
	}

	for _, pendingRole := range rolesToApply {
		role := pendingRole.role
		if role.Skip != true {
			if !role.PresentInRoles(member.Roles) {
				err := role.AddToCharacter(guild.Id, discordUserId, c.Discord.Session)
				if err != nil {
					return nil, fmt.Errorf("Error adding Discord role: %v", err)
				}
				text = append(text, fmt.Sprintf("__Adding role: **%s**__\n⮕ %s\n", role.Name, pendingRole.message))
			}
		}
	}

	if guild.SkipRemoval != true {
		for _, pendingRole := range rolesToRemove {
			role := pendingRole.role
			if role.Skip != true {
				if role.PresentInRoles(member.Roles) {
					err := role.RemoveFromCharacter(guild.Id, discordUserId, c.Discord.Session)
					if err != nil {
						return nil, fmt.Errorf("Error removing Discord role: %v", err)
					}
					text = append(text, fmt.Sprintf("__Removing role: **%s**__\n⮕ %s\n", role.Name, pendingRole.message))
				}
			}
		}
	}

	char.LastUpdateTime = time.Now()

	return text, nil
}
