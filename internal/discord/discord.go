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

func DMUser(s *discordgo.Session, i *discordgo.Interaction, message string) error {
	var userId string
	if i.Member != nil {
		userId = i.Member.User.ID
	}
	if i.User != nil {
		userId = i.User.ID
	}

	if userId == "" {
		return fmt.Errorf("Could not find user ID!")
	}
	channel, err := s.UserChannelCreate(userId)
	if err != nil {
		return fmt.Errorf("Could not create DM channel: %w", err)
	}
	_, err = s.ChannelMessageSend(channel.ID, message)
	if err != nil {
		return fmt.Errorf("Could not send message: %w", err)
	}
	return nil
}
