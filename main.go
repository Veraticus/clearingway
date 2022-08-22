package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"

	"github.com/Veraticus/clearingway/internal/discord"
	"github.com/Veraticus/clearingway/internal/fflogs"
	"github.com/Veraticus/clearingway/internal/ffxiv"

	"gopkg.in/yaml.v2"
)

func main() {
	discordToken, ok := os.LookupEnv("DISCORD_TOKEN")
	if !ok {
		panic("You must supply a DISCORD_TOKEN to start!")
	}

	discordChannelId, ok := os.LookupEnv("DISCORD_CHANNEL_ID")
	if !ok {
		panic("You must supply a DISCORD_CHANNEL_ID to start!")
	}

	fflogsClientId, ok := os.LookupEnv("FFLOGS_CLIENT_ID")
	if !ok {
		panic("You must supply a FFLOGS_CLIENT_ID to start!")
	}

	fflogsClientSecret, ok := os.LookupEnv("FFLOGS_CLIENT_SECRET")
	if !ok {
		panic("You must supply a FFLOGS_CLIENT_SECRET to start!")
	}

	fflogsInstance := fflogs.Init(fflogsClientId, fflogsClientSecret)

	encounters := &fflogs.Encounters{Encounters: []*fflogs.Encounter{}}
	config, err := ioutil.ReadFile("./config.yaml")
	if err != nil {
		panic(fmt.Errorf("Could not read config.yaml: %w", err))
	}
	yaml.Unmarshal(config, &encounters)
	encounters.Encounters = append(encounters.Encounters, fflogs.UltimateEncounters.Encounters...)
	for _, encounter := range encounters.Encounters {
		fmt.Printf("Encounter: %+v\n", encounter)
	}

	roles := &discord.Roles{Roles: []*discord.Role{}}
	roles.Roles = append(roles.Roles, discord.RolesForEncounters(encounters)...)
	roles.Roles = append(roles.Roles, discord.AllParsingRoles()...)
	roles.Roles = append(roles.Roles, discord.AllUltimateRoles()...)
	roles.Roles = append(roles.Roles, discord.AllServerRoles()...)

	d := &discord.Discord{
		Token:      discordToken,
		ChannelId:  discordChannelId,
		Fflogs:     fflogsInstance,
		Roles:      roles,
		Encounters: encounters,
		Characters: &ffxiv.Characters{Characters: map[string]*ffxiv.Character{}},
	}
	err = d.Start()
	defer d.Session.Close()
	if err != nil {
		panic(fmt.Errorf("Could not instantiate Discord: %w", err))
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}
