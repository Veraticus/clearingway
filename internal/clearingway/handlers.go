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
			fmt.Sprintf("Initialized in guild %s with no configuration!", gid)
			continue
		}
		existingRoles, err := s.GuildRoles(gid)
		if err != nil {
			fmt.Printf("Error getting existing roles: %v\n", err)
			return
		}

		fmt.Printf("Initializing roles...\n")
		guild.EncounterRoles = guild.Encounters.Roles()
		errs := guild.EncounterRoles.Ensure(gid, s, existingRoles)
		if len(errs) != 0 {
			fmt.Print("Error ensuring encounter roles:\n")
			for _, e := range errs {
				fmt.Printf("  %v", e)
			}
		}

		if guild.RelevantParsingEnabled {
			guild.RelevantParsingRoles = RelevantParsingRoles()
			errs = guild.RelevantParsingRoles.Ensure(gid, s, existingRoles)
			if len(errs) != 0 {
				fmt.Print("Error ensuring relevant parsing roles:\n")
				for _, e := range errs {
					fmt.Printf("  %v", e)
				}
			}
		}

		if guild.RelevantFlexingEnabled {
			guild.RelevantFlexingRoles = RelevantFlexingRoles()
			errs = guild.RelevantFlexingRoles.Ensure(gid, s, existingRoles)
			if len(errs) != 0 {
				fmt.Print("Error ensuring relevant flexing roles:\n")
				for _, e := range errs {
					fmt.Printf("  %v", e)
				}
			}
		}

		if guild.LegendEnabled {
			guild.LegendRoles = LegendRoles()
			errs = guild.LegendRoles.Ensure(gid, s, existingRoles)
			if len(errs) != 0 {
				fmt.Print("Error ensuring legend roles:\n")
				for _, e := range errs {
					fmt.Printf("  %v", e)
				}
			}
		}

		if guild.UltimateFlexingEnabled {
			guild.UltimateFlexingRoles = UltimateFlexingRoles()
			errs = guild.UltimateFlexingRoles.Ensure(gid, s, existingRoles)
			if len(errs) != 0 {
				fmt.Print("Error ensuring ultimate roles:\n")
				for _, e := range errs {
					fmt.Printf("  %v", e)
				}
			}
		}

		if guild.DatacenterEnabled {
			guild.DatacenterRoles = guild.Datacenters.AllRoles()
			errs = guild.DatacenterRoles.Ensure(gid, s, existingRoles)
			if len(errs) != 0 {
				fmt.Print("Error ensuring datacenter roles:\n")
				for _, e := range errs {
					fmt.Printf("  %v", e)
				}
			}
		}

		fmt.Printf("Adding commands...\n")
		_, err = s.ApplicationCommandCreate(event.User.ID, discordGuild.ID, ClearCommand)
		if err != nil {
			fmt.Printf("Could not add command: %v\n", err)
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

func (c *Clearingway) InteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		c.Clears(s, i)
	case discordgo.InteractionApplicationCommandAutocomplete:
		c.Autocomplete(s, i)
	}
}

func (c *Clearingway) Clears(s *discordgo.Session, i *discordgo.InteractionCreate) {
	g, ok := c.Guilds.Guilds[i.GuildID]
	if !ok {
		fmt.Printf("Interaction received from guild %s with no configuration!\n", i.GuildID)
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
			fmt.Printf("Error sending Discord message: %v", err)
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
		fmt.Printf("Error sending Discord message: %v", err)
		return
	}

	char, err := g.Characters.Init(world, firstName, lastName)
	if err != nil {
		err = discord.ContinueInteraction(s, i.Interaction, err.Error())
		if err != nil {
			fmt.Printf("Error sending Discord message: %v", err)
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
			fmt.Printf("Error sending Discord message: %v", err)
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
				fmt.Printf("Error sending Discord message: %v", err)
				return
			}
		}
	}

	err = discord.ContinueInteraction(s, i.Interaction,
		fmt.Sprintf("Verifying ownership of `%s (%s)`...", char.Name(), char.World),
	)
	if err != nil {
		fmt.Printf("Error sending Discord message: %v", err)
	}

	discordId := i.Member.User.ID
	isOwner, err := lodestone.CharacterIsOwnedByDiscordUser(char, discordId)
	if err != nil {
		err = discord.ContinueInteraction(s, i.Interaction, err.Error())
		if err != nil {
			fmt.Printf("Error sending Discord message: %v", err)
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
			fmt.Printf("Error sending Discord message: %v", err)
		}
		return
	}

	err = discord.ContinueInteraction(s, i.Interaction,
		fmt.Sprintf("Analyzing logs for `%s (%s)`...", char.Name(), char.World),
	)
	if err != nil {
		fmt.Printf("Error sending Discord message: %v", err)
	}

	if char.UpdatedRecently() {
		err = discord.ContinueInteraction(s, i.Interaction,
			fmt.Sprintf("Finished analysis for `%s (%s)`.", char.Name(), char.World),
		)
		if err != nil {
			fmt.Printf("Error sending Discord message: %v", err)
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

	err = discord.ContinueInteraction(s, i.Interaction,
		fmt.Sprintf("Finished analysis for `%s (%s)`.", char.Name(), char.World),
	)
	if err != nil {
		fmt.Printf("Error sending Discord message: %v", err)
	}

	for _, roleText := range roleTexts {
		err = discord.ContinueInteraction(s, i.Interaction, "_ _\n"+roleText)
		if err != nil {
			fmt.Printf("Error sending Discord message: %v", err)
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
		for _, world := range ffxiv.AllWorlds() {
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
		if !role.PresentInRoles(member.Roles) {
			err := role.AddToCharacter(guild.Id, discordUserId, c.Discord.Session, char)
			if err != nil {
				return nil, fmt.Errorf("Error adding Discord role: %v", err)
			}
			text = append(text, fmt.Sprintf("__Adding role: **%s**__\n⮕ %s\n", role.Name, pendingRole.message))
		}
	}

	for _, pendingRole := range rolesToRemove {
		role := pendingRole.role
		if role.PresentInRoles(member.Roles) {
			err := role.RemoveFromCharacter(guild.Id, discordUserId, c.Discord.Session, char)
			if err != nil {
				return nil, fmt.Errorf("Error removing Discord role: %v", err)
			}
			text = append(text, fmt.Sprintf("__Removing role: **%s**__\n⮕ %s\n", role.Name, pendingRole.message))
		}
	}

	char.LastUpdateTime = time.Now()

	return text, nil
}
