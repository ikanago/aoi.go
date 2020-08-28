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

// TweetCreate represents parsed results of `twitter create` command.
type TweetCreate struct {
	ID       string
	Keywords []string
}

// TweetAdd represents parsed results of `twitter add` command.
type TweetAdd struct {
	ID       string
	Keywords []string
}

// TweetRemove represents parsed results of `twitter remove` command.
type TweetRemove struct {
	ID       string
	Keywords []string
}

// TweetDelete represents parsed results of `twitter delete` command.
type TweetDelete struct {
	ID string
}

// TweetChange represents parsed results of `twitter change` command.
type TweetChange struct {
	ID      string
	Channel string
}

// TweetShow represents parsed results of `twitter show` command.
type TweetShow struct{}

func (Help) handle(session *discordgo.Session, message *discordgo.Message) (err error) {
	messageEmbed := discordgo.MessageEmbed{
		Color:  0x4bede7,
		Type:   discordgo.EmbedTypeRich,
		Title:  "アオイチャンのコマンド",
		Fields: helpMessageEmbeds,
	}
	_, err = session.ChannelMessageSendEmbed(message.ChannelID, &messageEmbed)
	return
}

func (Ping) handle(session *discordgo.Session, message *discordgo.Message) (err error) {
	_, err = session.ChannelMessageSend(message.ChannelID, "Pong!")
	return
}

func (tweetCreate TweetCreate) handle(session *discordgo.Session, message *discordgo.Message) (err error) {
	return
}

func (tweetAdd TweetAdd) handle(session *discordgo.Session, message *discordgo.Message) (err error) {
	return
}

func (tweetRemove TweetRemove) handle(session *discordgo.Session, message *discordgo.Message) (err error) {
	return
}

func (tweetDelete TweetDelete) handle(session *discordgo.Session, message *discordgo.Message) (err error) {
	return
}

func (tweetChange TweetChange) handle(session *discordgo.Session, message *discordgo.Message) (err error) {
	return
}

func (tweetShow TweetShow) handle(session *discordgo.Session, message *discordgo.Message) (err error) {
	return
}
