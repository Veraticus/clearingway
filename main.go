package main

import (
	"fmt"
	"os"

	"github.com/Veraticus/clearingway/internal/discord"
	"github.com/Veraticus/clearingway/internal/fflogs"
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

	encounters, ok := os.LookupEnv("ENCOUNTERS")
	if !ok {
		panic("You must supply ENCOUNTERS to start!")
	}

	discord := &discord.Discord{
		Token: discordToken,
	}
	err := discord.Start()
	defer discord.Session.Close()
	if err != nil {
		panic(fmt.Errorf("Could not instantiate Discord: %f", err))
	}

	fflogs := fflogs.Init(fflogsClientId, fflogsClientSecret)
	fmt.Printf("Fflogs is %+v", fflogs)
	fmt.Printf("Encounters is %+v", encounters)
}
