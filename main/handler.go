package main

import (
	"errors"
	"fmt"
	"strings"

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
	ScreenName string
	Keywords   []string
}

// TweetAdd represents parsed results of `twitter add` command.
type TweetAdd struct {
	ScreenName string
	Keywords   []string
}

// TweetRemove represents parsed results of `twitter remove` command.
type TweetRemove struct {
	ScreenName string
	Keywords   []string
}

// TweetDelete represents parsed results of `twitter delete` command.
type TweetDelete struct {
	ScreenName string
}

// TweetChange represents parsed results of `twitter change` command.
type TweetChange struct {
	ScreenName string
	Channel    string
}

// TweetShow represents parsed results of `twitter show` command.
type TweetShow struct{}

// MemoRegister represents parsed results of `memo TEXT` command.
type MemoRegister struct {
	Text string
}

// MemoShow represents parsed results of `memo TEXT` command.
type MemoShow struct{}

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
	err = createFilter(tweetCreate.ScreenName, tweetCreate.Keywords, message.ChannelID)
	if err != nil {
		return
	}
	reply := fmt.Sprintf("@%s のフィルターを作成しました\n現在のキーワード: %s", tweetCreate.ScreenName, strings.Join(tweetCreate.Keywords, ", "))
	_, err = session.ChannelMessageSend(message.ChannelID, reply)
	return
}

func (tweetAdd TweetAdd) handle(session *discordgo.Session, message *discordgo.Message) (err error) {
	updatedKeywords, err := addFilter(tweetAdd.ScreenName, tweetAdd.Keywords)
	if err != nil {
		return
	}
	reply := fmt.Sprintf("@%s のフィルターを更新しました\n現在のキーワード: %s", tweetAdd.ScreenName, strings.Join(updatedKeywords, ", "))
	_, err = session.ChannelMessageSend(message.ChannelID, reply)
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

func (memoShow MemoShow) handle(session *discordgo.Session, message *discordgo.Message) (err error) {
	memos, err := fetchMemo(message.ChannelID)
	if err != nil {
		return
	}
	if len(memos) == 0 {
		return errors.New("このチャンネルにはまだ発言が記録されていません")
	}

	var reply string
	for i, memo := range memos {
		reply += fmt.Sprintf("%2d: %s (%s)\n", i, memo.Text, memo.Timestamp.Format("2006-01-02"))
	}
	_, err = session.ChannelMessageSend(message.ChannelID, reply)
	return
}

func (memoRegister MemoRegister) handle(session *discordgo.Session, message *discordgo.Message) (err error) {
	reply, err := addMemo(message.ChannelID, memoRegister.Text)
	if err != nil {
		return
	}

	_, err = session.ChannelMessageSend(message.ChannelID, reply)
	return
}
