package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

var twitterClient *twitter.Client

func initClient(credential *Credential) error {
	if credential.ConsumerKey == "" || credential.ConsumerSecret == "" || credential.AccessToken == "" || credential.AccessTokenSecret == "" {
		return errors.New("Twitter API tokens not specified")
	}

	config := oauth1.NewConfig(credential.ConsumerKey, credential.ConsumerSecret)
	token := oauth1.NewToken(credential.AccessToken, credential.AccessTokenSecret)
	httpClient := config.Client(oauth1.NoContext, token)
	twitterClient = twitter.NewClient(httpClient)
	return nil
}

func initStream(session *discordgo.Session) (stream *twitter.Stream, demux twitter.SwitchDemux, err error) {
	demux = twitter.NewSwitchDemux()
	demux.Tweet = func(tweet *twitter.Tweet) {
		isSatisfy, channelID := filterTweet(tweet.User.ScreenName, tweet.Text)
		if isSatisfy {
			tweetURL := fmt.Sprintf("https://twitter.com/%s/status/%s", tweet.User.ScreenName, tweet.IDStr)
			reply := tweet.Text + "\n" + tweetURL
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

func getUserID(screenName string) (id string, err error) {
	users, _, err := twitterClient.Users.Lookup(&twitter.UserLookupParams{
		ScreenName: []string{screenName},
	})
	if err != nil {
		return
	}
	id = users[0].IDStr
	return
}

func filterTweet(screenName string, tweet string) (isSatisfy bool, channelID string) {
	filter := GetFilter(screenName)
	if filter == nil {
		return false, ""
	}

	for _, keyword := range filter.Keywords {
		if strings.Contains(tweet, keyword) {
			return true, filter.ChannelID
		}
	}
	return false, ""
}
