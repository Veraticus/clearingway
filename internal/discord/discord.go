package discord

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type Discord struct {
	Token   string
	Session *discordgo.Session
}

func (d *Discord) Start() error {
	s, err := discordgo.New("Bot " + d.Token)
	if err != nil {
		return fmt.Errorf("Could not start Discord: %f", err)
	}

	s.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages

	d.Session = s
	return nil
}

func StartInteraction(s *discordgo.Session, i *discordgo.Interaction, message string) error {
	err := s.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message,
		},
	})
	return err
}

func ContinueInteraction(s *discordgo.Session, i *discordgo.Interaction, message string) error {
	_, err := s.FollowupMessageCreate(i, true, &discordgo.WebhookParams{
		Content: message,
	})
	return err
}
