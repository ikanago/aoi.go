package main

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

var twitterClient *twitter.Client

// InitClient initializes Twitter Client.
func InitClient(credential *Credential) error {
	if credential.ConsumerKey == "" || credential.ConsumerSecret == "" || credential.AccessToken == "" || credential.AccessTokenSecret == "" {
		return errors.New("Twitter API tokens not specified")
	}

	config := oauth1.NewConfig(credential.ConsumerKey, credential.ConsumerSecret)
	token := oauth1.NewToken(credential.AccessToken, credential.AccessTokenSecret)
	httpClient := config.Client(oauth1.NoContext, token)
	twitterClient = twitter.NewClient(httpClient)
	return nil
}

// InitStream initializes tweet stream.
func InitStream(session *discordgo.Session) (stream *twitter.Stream, demux twitter.SwitchDemux, err error) {
	demux = twitter.NewSwitchDemux()
	demux.Tweet = func(tweet *twitter.Tweet) {
		filter := GetFilter(tweet.User.ScreenName)
		isSatisfy, channelID := filterTweet(filter, tweet.Text)
		if isSatisfy {
			reply := formatTweet(tweet.Text, tweet.User.ScreenName, tweet.IDStr)
			_, err = session.ChannelMessageSend(channelID, reply)
			return
		}
	}

	ids := GetAllIDs()
	filterParam := &twitter.StreamFilterParams{
		Follow: ids,
	}
	stream, err = twitterClient.Streams.Filter(filterParam)
	return
}

func filterTweet(filter *FilterDocument, tweet string) (isSatisfy bool, channelID string) {
	if filter == nil {
		return false, ""
	}

	// Exclude retweets
	if strings.HasPrefix(tweet, "RT") {
		return false, ""
	}

	for _, keyword := range filter.Keywords {
		if strings.Contains(tweet, keyword) {
			return true, filter.ChannelID
		}
	}
	return false, ""
}

const hyperLinkPattern = `(http|https):\/\/([a-zA-Z0-9.]/?)+`

var hyperLinkRegExp = regexp.MustCompile(hyperLinkPattern)

func formatTweet(text, screenName, id string) string {
	text = string(hyperLinkRegExp.ReplaceAll([]byte(text), []byte("")))
	return fmt.Sprintf("%s\nhttps://twitter.com/%s/status/%s", text, screenName, id)
}

// GetUserID queries id string of a specific screen name.
func GetUserID(screenName string) (id string, err error) {
	users, _, err := twitterClient.Users.Lookup(&twitter.UserLookupParams{
		ScreenName: []string{screenName},
	})
	if err != nil {
		return
	}
	id = users[0].IDStr
	return
}
