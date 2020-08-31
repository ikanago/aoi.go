package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// Command is an interface of parsed commands.
type Command interface {
	Handle(*discordgo.Session, *discordgo.Message) error
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

// Handle deals `help` command.
func (Help) Handle(session *discordgo.Session, message *discordgo.Message) (err error) {
	messageEmbed := discordgo.MessageEmbed{
		Color:  0x4bede7,
		Type:   discordgo.EmbedTypeRich,
		Title:  "アオイチャンのコマンド",
		Fields: HelpMessageEmbeds,
	}
	_, err = session.ChannelMessageSendEmbed(message.ChannelID, &messageEmbed)
	return
}

// Handle deals `ping` command.
func (Ping) Handle(session *discordgo.Session, message *discordgo.Message) (err error) {
	_, err = session.ChannelMessageSend(message.ChannelID, "Pong!")
	return
}

// Handle deals `tweet create` command.
func (tweetCreate TweetCreate) Handle(session *discordgo.Session, message *discordgo.Message) (err error) {
	err = CreateFilter(tweetCreate.ScreenName, tweetCreate.Keywords, message.ChannelID)
	if err != nil {
		return
	}
	reply := fmt.Sprintf("@%s のフィルターを作成しました\n現在のキーワード: %s", tweetCreate.ScreenName, strings.Join(tweetCreate.Keywords, ", "))
	_, err = session.ChannelMessageSend(message.ChannelID, reply)
	return
}

// Handle deals `tweet add` command.
func (tweetAdd TweetAdd) Handle(session *discordgo.Session, message *discordgo.Message) (err error) {
	updatedKeywords, err := AddFilter(tweetAdd.ScreenName, tweetAdd.Keywords)
	if err != nil {
		return
	}
	reply := fmt.Sprintf("@%s のフィルターを更新しました\n現在のキーワード: %s", tweetAdd.ScreenName, strings.Join(updatedKeywords, ", "))
	_, err = session.ChannelMessageSend(message.ChannelID, reply)
	return
}

// Handle deals `tweet remove` command.
func (tweetRemove TweetRemove) Handle(session *discordgo.Session, message *discordgo.Message) (err error) {
	return
}

// Handle deals `tweet delete` command.
func (tweetDelete TweetDelete) Handle(session *discordgo.Session, message *discordgo.Message) (err error) {
	return
}

// Handle deals `tweet change` command.
func (tweetChange TweetChange) Handle(session *discordgo.Session, message *discordgo.Message) (err error) {
	return
}

// Handle deals `tweet show` command.
func (tweetShow TweetShow) Handle(session *discordgo.Session, message *discordgo.Message) (err error) {
	return
}

// Handle deals `memo show` command.
func (memoShow MemoShow) Handle(session *discordgo.Session, message *discordgo.Message) (err error) {
	memos, err := FetchMemo(message.ChannelID)
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

// Handle deals `memo` command.
func (memoRegister MemoRegister) Handle(session *discordgo.Session, message *discordgo.Message) (err error) {
	reply, err := CreateMemo(message.ChannelID, memoRegister.Text)
	if err != nil {
		return
	}

	_, err = session.ChannelMessageSend(message.ChannelID, reply)
	return
}
