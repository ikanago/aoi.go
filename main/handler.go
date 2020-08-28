package main

import (
	"github.com/bwmarrin/discordgo"
)

// Command is an interface of parsed commands.
type Command interface {
	handle(*discordgo.Session, *discordgo.Message) error
}

// Help represents parsed results of `help` command.
type Help struct{}

// Ping represents parsed results of `ping` command.
type Ping struct{}

func (Help) handle(session *discordgo.Session, message *discordgo.Message) (err error) {
	_, err = session.ChannelMessageSend(message.ChannelID, "help")
	return
}

func (Ping) handle(session *discordgo.Session, message *discordgo.Message) (err error) {
	_, err = session.ChannelMessageSend(message.ChannelID, "Pong!")
	return
}
