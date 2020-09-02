package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFilterTweet(t *testing.T) {
	t.Run("Invalid filter", func(t *testing.T) {
		tweet := ""
		expectedBool, expectedString := false, ""
		actualBool, actualString := filterTweet(nil, tweet)
		assert := assert.New(t)
		assert.Equal(expectedBool, actualBool)
		assert.Equal(expectedString, actualString)
	})

	t.Run("Matched", func(t *testing.T) {
		filter := &FilterDocument{
			ID:         "1111",
			ScreenName: "xxxx",
			Keywords:   []string{"創作2コマ漫画", "ねこ"},
			ChannelID:  "asdfghjkl",
		}
		tweet := "創作2コマ漫画　そのxxx https://example.com"
		expectedBool, expectedString := true, "asdfghjkl"
		actualBool, actualString := filterTweet(filter, tweet)
		assert := assert.New(t)
		assert.Equal(expectedBool, actualBool)
		assert.Equal(expectedString, actualString)
	})

	t.Run("Not matched", func(t *testing.T) {
		filter := &FilterDocument{
			ID:         "1111",
			ScreenName: "xxxx",
			Keywords:   []string{"創作2コマ漫画", "ねこ"},
			ChannelID:  "asdfghjkl",
		}
		tweet := "RT @70_pocky: 創作2コマ漫画　そのxxx https://example.com"
		expectedBool, expectedString := false, ""
		actualBool, actualString := filterTweet(filter, tweet)
		assert := assert.New(t)
		assert.Equal(expectedBool, actualBool)
		assert.Equal(expectedString, actualString)
	})
}
