package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Veraticus/clearingway/internal/clearingway"
	"github.com/Veraticus/clearingway/internal/discord"
	"github.com/Veraticus/clearingway/internal/fflogs"
	"github.com/Veraticus/clearingway/internal/lodestone"

	"gopkg.in/yaml.v3"
)

func main() {
	discordToken, ok := os.LookupEnv("DISCORD_TOKEN")
	if !ok {
		panic("You must supply a DISCORD_TOKEN to start!")
	}

	fflogsClientId, ok := os.LookupEnv("FFLOGS_CLIENT_ID")
	if !ok {
		panic("You must supply a FFLOGS_CLIENT_ID to start!")
	}

	fflogsClientSecret, ok := os.LookupEnv("FFLOGS_CLIENT_SECRET")
	if !ok {
		panic("You must supply a FFLOGS_CLIENT_SECRET to start!")
	}

	c := &clearingway.Clearingway{
		Config: &clearingway.Config{},
		Fflogs: fflogs.Init(fflogsClientId, fflogsClientSecret),
		Discord: &discord.Discord{
			Token: discordToken,
		},
	}

	config, err := ioutil.ReadFile("./config.yaml")
	if err != nil {
		panic(fmt.Errorf("Could not read config.yaml: %w", err))
	}
	err = yaml.Unmarshal(config, &c.Config)
	if err != nil {
		panic(fmt.Errorf("Could not unmarshal config.yaml: %w", err))
	}

	c.Init()

	fmt.Printf("Clearingway is: %+v\n", c)
	for _, guild := range c.Guilds.Guilds {
		fmt.Printf("Guild added: %+v\n", guild)

		if guild.EncounterRoles != nil {
			fmt.Printf("Encounter roles: %+v\n", guild.EncounterRoles.Roles)
		}

		if guild.RelevantParsingRoles != nil {
			fmt.Printf("Relevant parsing roles: %+v\n", guild.RelevantParsingRoles.Roles)
		}

		if guild.RelevantFlexingRoles != nil {
			fmt.Printf("Relevant flexing roles: %+v\n", guild.RelevantFlexingRoles.Roles)
		}

		if guild.LegendRoles != nil {
			fmt.Printf("Legend roles: %+v\n", guild.LegendRoles.Roles)
		}

		if guild.UltimateFlexingRoles != nil {
			fmt.Printf("Ultimate flexing roles: %+v\n", guild.UltimateFlexingRoles.Roles)
		}

		if guild.DatacenterRoles != nil {
			fmt.Printf("Datacenter roles: %+v\n", guild.DatacenterRoles.Roles)
		}
	}

	fmt.Printf("Starting Discord...\n")
	err = c.Discord.Start()
	if err != nil {
		panic(fmt.Errorf("Could not instantiate Discord: %w", err))
	}
	defer c.Discord.Session.Close()

	var arg string
	args := os.Args[1:]
	if len(args) == 0 {
		arg = ""
	} else {
		arg = args[0]
	}
	switch arg {
	case "run":
		run(c)
	default:
		start(c)
	}
}

func start(c *clearingway.Clearingway) {
	c.Discord.Session.AddHandler(c.DiscordReady)
	c.Discord.Session.AddHandler(c.InteractionCreate)
	err := c.Discord.Session.Open()
	if err != nil {
		panic(fmt.Errorf("Could not open Discord session: %f", err))
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}

func run(c *clearingway.Clearingway) {
	if len(os.Args) != 7 {
		panic("Provide a world, firstName, lastName, guildId, and discordId!")
	}
	world := os.Args[2]
	firstName := os.Args[3]
	lastName := os.Args[4]
	guildId := os.Args[5]
	discordId := os.Args[6]

	guild, ok := c.Guilds.Guilds[guildId]
	if !ok {
		panic(fmt.Sprintf("Guild %s not setup in config.yaml but you tried to run me in it!", guildId))
	}

	c.Discord.Session.AddHandler(c.DiscordReady)
	err := c.Discord.Session.Open()
	if err != nil {
		panic(fmt.Errorf("Could not open Discord session: %f", err))
	}

	for c.Ready != true {
		fmt.Printf("Waiting for Clearingway to be ready...\n")
		time.Sleep(2 * time.Second)
	}

	char, err := guild.Characters.Init(world, firstName, lastName)
	if err != nil {
		panic(err)
	}

	err = c.Fflogs.SetCharacterLodestoneID(char)
	if err != nil {
		fmt.Printf("Could not find character in FF Logs: %+v\n", err)
		err = lodestone.SetCharacterLodestoneID(char)
		if err != nil {
			panic(fmt.Errorf("Could not find character in the Lodestone: %+v", err))
		}
	}

	isOwner, err := lodestone.CharacterIsOwnedByDiscordUser(char, discordId)
	if err != nil {
		panic(err)
	}
	if !isOwner {
		panic("That character is not owned by that Discord ID!")
	}

	roleTexts, err := c.UpdateCharacterInGuild(char, discordId, guild)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Character %s (%s) updated in guild %s.\n", char.Name(), char.World, guild.Name)

	for _, roleText := range roleTexts {
		fmt.Printf(roleText + "\n")
	}
}
