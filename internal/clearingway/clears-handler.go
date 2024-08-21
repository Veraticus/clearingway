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

	err := discord.StartInteraction(s, i.Interaction, "Received `/clears`...")
	if err != nil {
		fmt.Printf("Error sending Discord message: %v\n", err)
		return
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

	c.ClearsHelper(s, i, g, world, firstName, lastName)
}

func (c *Clearingway) ClearsHelper(s *discordgo.Session, i *discordgo.InteractionCreate, g *Guild, world string, firstName string, lastName string) {
	if len(world) == 0 || len(firstName) == 0 || len(lastName) == 0 {
		err := discord.ContinueInteraction(s, i.Interaction, "`/clears` command failed! Please input your world, first name, and last name.")
		if err != nil {
			fmt.Printf("Error sending Discord message: %v\n", err)
		}
		return
	}

	title := cases.Title(language.AmericanEnglish)
	world = title.String(cleanWorld(world))
	firstName = title.String(firstName)
	lastName = title.String(lastName)

	if !ffxiv.IsWorld(world) {
		err := discord.ContinueInteraction(s, i.Interaction,
			fmt.Sprintf("`%s` is not a valid world! Make sure you spelled your world name properly.", world),
		)
		if err != nil {
			fmt.Printf("Error sending Discord message: %v\n", err)
		}
		return
	}

	err := discord.ContinueInteraction(s, i.Interaction,
		fmt.Sprintf("Finding `%s %s (%s)` in the Lodestone...", firstName, lastName, world),
	)
	if err != nil {
		fmt.Printf("Error sending Discord message: %v\n", err)
		return
	}

	char, err := g.Characters.Init(world, firstName, lastName)
	if err != nil {
		err := discord.ContinueInteraction(s, i.Interaction, err.Error())
		if err != nil {
			fmt.Printf("Error sending Discord message: %v\n", err)
		}
		return
	}

	err = c.Fflogs.SetCharacterLodestoneID(char)
	if err != nil {
		err := discord.ContinueInteraction(s, i.Interaction,
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
			err := discord.ContinueInteraction(s, i.Interaction,
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
		err = discord.ContinueInteraction(s, i.Interaction,
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

	roleTexts, err := c.UpdateClearsForCharacterInGuild(char, i.Member.User.ID, g)
	if err != nil {
		err = discord.ContinueInteraction(s, i.Interaction,
			fmt.Sprintf("Could not analyze clears for `%s (%s)`: %s", char.Name(), char.World, err),
		)
		if err != nil {
			fmt.Printf("Error sending Discord message: %v\n", err)
		}
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
}

func (c *Clearingway) UpdateClearsForCharacterInGuild(
	char *ffxiv.Character,
	discordUserId string,
	guild *Guild,
) ([]string, error) {
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
		fmt.Printf("%s (%d):\n", e.Name, len(e.Ranks(rankings)))
		for _, r := range e.Ranks(rankings) {
			fmt.Printf("  %+v\n", r)
		}
	}

	fmt.Printf("Found the following ultimate rankings for %s (%s)...\n", char.Name(), char.World)
	for _, e := range UltimateEncounters.Encounters {
		fmt.Printf("%s (%d):\n", e.Name, len(e.Ranks(rankings)))
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

	// Add achievement roles when we have a performant way to do it
	// fmt.Printf("Scraping completed achievements...\n")
	// clearedAchievements, err := lodestone.GetAchievements(char)
	// if err != nil {
	// 	fmt.Printf("Please check to make sure your achievements page is not private.")
	// } else {
	// 	for _, role := range guild.AchievementRoles.Roles {
	// 		for _, clearedAchievement := range clearedAchievements {
	// 			if clearedAchievement == role.Description {
	// 				message := fmt.Sprintf("Cleared achievement %s", clearedAchievement)
	// 				text = append(text, fmt.Sprintf("__Adding role: **%s**__\n-> %s\n", role.Name, message))
	// 				rolesToApply = append(rolesToApply, &pendingRole{role: role, message: message})
	// 			}
	// 		}
	// 	}
	// }
	fmt.Printf("Scraping completed.\n")

	for _, pendingRole := range rolesToApply {
		role := pendingRole.role
		if !role.Skip {
			if !role.PresentInRoles(member.Roles) {
				err := role.AddToCharacter(guild.Id, discordUserId, c.Discord.Session)
				if err != nil {
					return nil, fmt.Errorf("Error adding Discord role +%v: %w", role, err)
				}
				text = append(text, fmt.Sprintf("__Adding role: **%s**__\n⮕ %s\n", role.Name, pendingRole.message))
			}
		}
	}

	if !guild.SkipRemoval {
		for _, pendingRole := range rolesToRemove {
			role := pendingRole.role
			if !role.Skip {
				if role.PresentInRoles(member.Roles) {
					err := role.RemoveFromCharacter(guild.Id, discordUserId, c.Discord.Session)
					if err != nil {
						return nil, fmt.Errorf("Error adding Discord role +%v: %w", role, err)
					}
					text = append(text, fmt.Sprintf("__Removing role: **%s**__\n⮕ %s\n", role.Name, pendingRole.message))
				}
			}
		}
	}

	char.LastUpdateTime = time.Now()

	return text, nil
}
