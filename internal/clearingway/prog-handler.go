package clearingway

import (
	"fmt"
	"strings"
	"time"

	"github.com/Veraticus/clearingway/internal/discord"
	"github.com/Veraticus/clearingway/internal/fflogs"
	"github.com/Veraticus/clearingway/internal/ffxiv"
	"github.com/Veraticus/clearingway/internal/lodestone"

	"github.com/bwmarrin/discordgo"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func (c *Clearingway) Prog(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
	var reportId string

	if option, ok := optionMap["world"]; ok {
		world = option.StringValue()
	}
	if option, ok := optionMap["first-name"]; ok {
		firstName = option.StringValue()
	}
	if option, ok := optionMap["last-name"]; ok {
		lastName = option.StringValue()
	}
	if option, ok := optionMap["report-id"]; ok {
		reportId = option.StringValue()
	}

	if len(world) == 0 || len(firstName) == 0 || len(lastName) == 0 || len(reportId) == 0 {
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "`/prog` command failed! Make sure you input your world, first name, last name, and fflogs report URL or ID.",
			},
		})
		if err != nil {
			fmt.Printf("Error sending Discord message: %v\n", err)
		}
		return
	}

	title := cases.Title(language.AmericanEnglish)
	world = title.String(cleanWorld(world))
	firstName = title.String(firstName)
	lastName = title.String(lastName)
	reportId = CleanReportId(reportId)

	if !ffxiv.IsWorld(world) {
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("`%s` is not a valid world! Make sure you spelled your world name properly.", world),
			},
		})
		if err != nil {
			fmt.Printf("Error sending Discord message: %v\n", err)
		}
		return
	}

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
				"I could not verify your ownership of `%s (%s)`!\nIf this is your character, add the following code to your Lodestone profile and then run `/prog` again:\n\n**%s**\n\nYou can edit your Lodestone profile at https://na.finalfantasyxiv.com/lodestone/my/setting/profile/",
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
		fmt.Sprintf("Analyzing report for `%s (%s)`...", char.Name(), char.World),
	)
	if err != nil {
		fmt.Printf("Error sending Discord message: %v\n", err)
	}

	if char.UpdatedRecently() {
		err = discord.ContinueInteraction(s, i.Interaction,
			fmt.Sprintf("Finished report analysis for `%s (%s)`.", char.Name(), char.World),
		)
		if err != nil {
			fmt.Printf("Error sending Discord message: %v\n", err)
		}
		return
	}

	roleTexts, err := c.UpdateProgForCharacterInGuild(reportId, char, i.Member.User.ID, g)
	if err != nil {
		err = discord.ContinueInteraction(s, i.Interaction,
			fmt.Sprintf("Could not analyze prog for `%s (%s)`: %s", char.Name(), char.World, err),
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

func (c *Clearingway) UpdateProgForCharacterInGuild(
	reportId string,
	char *ffxiv.Character,
	discordUserId string,
	guild *Guild,
) ([]string, error) {
	rankingsToGet := []*fflogs.RankingToGet{}
	for _, encounter := range guild.AllEncounters() {
		rankingsToGet = append(rankingsToGet, &fflogs.RankingToGet{IDs: encounter.Ids, Difficulty: encounter.DifficultyInt()})
	}
	fights, err := c.Fflogs.GetProgForReport(reportId, rankingsToGet, char)
	if err != nil {
		return nil, fmt.Errorf("Error retrieving prog: %w", err)
	}

	fmt.Printf("Found the following relevant fights for %s (%s)...\n", char.Name(), char.World)
	for _, e := range guild.Encounters.Encounters {
		for _, r := range e.Fights(fights) {
			fmt.Printf("  %+v\n", r)
		}
	}

	member, err := c.Discord.Session.GuildMember(guild.Id, discordUserId)
	if err != nil {
		return nil, fmt.Errorf("Could not retrieve roles for user: %w", err)
	}
	existingRoles := &Roles{Roles: []*Role{}}
	for _, guildRole := range guild.AllRoles() {
		if guildRole.Skip {
			continue
		}
		for _, memberDiscordRole := range member.Roles {
			if memberDiscordRole == guildRole.DiscordRole.ID {
				existingRoles.Roles = append(existingRoles.Roles, guildRole)
			}
		}
	}
	text := []string{}

	shouldApplyOpts := &ShouldApplyOpts{
		Character:     char,
		Fights:        fights,
		ExistingRoles: existingRoles,
		Encounters:    guild.Encounters,
	}

	for _, encounter := range guild.Encounters.Encounters {
		if encounter.ProgRoles == nil {
			continue
		}

		shouldApply, message, rolesToApply, rolesToRemove := encounter.ProgRoles.ShouldApply(shouldApplyOpts)

		text = append(text, message)

		if shouldApply {
			for _, role := range rolesToApply {
				if role.Skip != true {
					if !role.PresentInRoles(member.Roles) {
						err := role.AddToCharacter(guild.Id, discordUserId, c.Discord.Session)
						if err != nil {
							return nil, fmt.Errorf("Error adding Discord role: %v", err)
						}
						text = append(text, fmt.Sprintf("Adding role: __**%s**__\n", role.Name))
					}
				}
			}

			for _, role := range rolesToRemove {
				if role.Skip != true {
					if role.PresentInRoles(member.Roles) {
						err := role.RemoveFromCharacter(guild.Id, discordUserId, c.Discord.Session)
						if err != nil {
							return nil, fmt.Errorf("Error removing Discord role: %v", err)
						}
						text = append(text, fmt.Sprintf("Removing role: __**%s**__\n", role.Name))
					}
				}
			}
		}
	}

	char.LastUpdateTime = time.Now()

	return text, nil
}

func CleanReportId(reportId string) string {
	reportId = strings.TrimRight(reportId, "/")
	reportIds := strings.Split(reportId, "#")
	reportId = reportIds[0]
	reportIds = strings.Split(reportId, "/")
	reportId = reportIds[len(reportIds)-1]
	reportIds = strings.Split(reportId, "#")
	reportId = reportIds[0]
	return reportId
}
