package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/Veraticus/clearingway/internal/discord"
	"github.com/Veraticus/clearingway/internal/fflogs"
	"github.com/Veraticus/clearingway/internal/ffxiv"
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

	// Encounters should be in `<name>=<encounterId>` format separated by commas.
	encounters, ok := os.LookupEnv("ENCOUNTERS")
	if !ok {
		panic("You must supply ENCOUNTERS to start!")
	}
	relevantEncounters := &fflogs.Encounters{Encounters: []*fflogs.Encounter{}}
	for _, encounter := range strings.Split(encounters, ",") {
		e := strings.Split(encounter, "=")
		name := e[0]
		id, err := strconv.Atoi(e[1])
		if err != nil {
			panic(fmt.Errorf("Could not convert encounter to integer for %v: %w", name, err))
		}
		relevantEncounters.Encounters = append(relevantEncounters.Encounters, &fflogs.Encounter{Name: name, IDs: []int{id}})
	}
	for _, relevantEncounter := range relevantEncounters.Encounters {
		fmt.Printf("Relevant encounter: %+v\n", relevantEncounter)
	}

	roles := &discord.Roles{Roles: []*discord.Role{}}
	roles.Roles = append(roles.Roles, discord.AllParsingRoles()...)
	roles.Roles = append(roles.Roles, discord.AllUltimateRoles()...)
	roles.Roles = append(roles.Roles, discord.RolesForEncounters(relevantEncounters)...)
	roles.Roles = append(roles.Roles, discord.AllServerRoles()...)

	discord := &discord.Discord{
		Token:              discordToken,
		ChannelId:          discordChannelId,
		RelevantEncounters: relevantEncounters,
		Fflogs:             fflogsInstance,
		Roles:              roles,
		Characters:         &ffxiv.Characters{Characters: map[string]*ffxiv.Character{}},
	}
	err := discord.Start()
	defer discord.Session.Close()
	if err != nil {
		panic(fmt.Errorf("Could not instantiate Discord: %w", err))
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}
