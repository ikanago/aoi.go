package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseBasicCommand(t *testing.T) {
	t.Run("Command not specified", func(t *testing.T) {
		input := "<@!1234567890>"
		actual, err := ParseCommand(input)
		assert := assert.New(t)
		assert.Nil(actual)
		assert.NotNil(err)
	})

	t.Run("Help command", func(t *testing.T) {
		input := "<@!1234567890> help"
		expected := Help{}
		actual, err := ParseCommand(input)
		assert := assert.New(t)
		assert.Equal(expected, actual)
		assert.Nil(err)
	})

	t.Run("Command not specified", func(t *testing.T) {
		input := "<@!1234567890> ping"
		expected := Ping{}
		actual, err := ParseCommand(input)
		assert := assert.New(t)
		assert.Equal(expected, actual)
		assert.Nil(err)
	})

	t.Run("Unknown command", func(t *testing.T) {
		input := "<@!1234567890> hoge"
		actual, err := ParseCommand(input)
		assert := assert.New(t)
		assert.Nil(actual)
		assert.NotNil(err)
	})
}

func TestParseTweetCommand(t *testing.T) {
	t.Run("Subcommand not specified", func(t *testing.T) {
		input := []string{"<@!1234567890>", "tweet"}
		actual, err := parseTweetCommand(input)
		assert := assert.New(t)
		assert.Nil(actual)
		assert.NotNil(err)
	})

	t.Run("Create command", func(t *testing.T) {
		input := []string{"<@!1234567890>", "tweet", "create", "ab_12_cd", "地震", "Aoi"}
		expected := TweetCreate{
			ID: "ab_12_cd",
			Keywords: []string{"地震", "Aoi"},
		}
		actual, err := parseTweetCommand(input)
		assert := assert.New(t)
		assert.Equal(expected, actual)
		assert.Nil(err)
	})

	t.Run("Add command", func(t *testing.T) {
		input := []string{"<@!1234567890>", "tweet", "add", "docker", "container", "k8s", "cloud"}
		expected := TweetAdd{
			ID: "docker",
			Keywords: []string{"container", "k8s", "cloud"},
		}
		actual, err := parseTweetCommand(input)
		assert := assert.New(t)
		assert.Equal(expected, actual)
		assert.Nil(err)
	})

	t.Run("Remove command", func(t *testing.T) {
		input := []string{"<@!1234567890>", "tweet", "remove", "docker", "container", "k8s"}
		expected := TweetRemove{
			ID: "docker",
			Keywords: []string{"container", "k8s"},
		}
		actual, err := parseTweetCommand(input)
		assert := assert.New(t)
		assert.Equal(expected, actual)
		assert.Nil(err)
	})

	t.Run("Delete command", func(t *testing.T) {
		input := []string{"<@!1234567890>", "tweet", "delete", "docker"}
		expected := TweetDelete{
			ID: "docker",
		}
		actual, err := parseTweetCommand(input)
		assert := assert.New(t)
		assert.Equal(expected, actual)
		assert.Nil(err)
	})

	t.Run("Change command", func(t *testing.T) {
		input := []string{"<@!1234567890>", "tweet", "change", "docker", "#container"}
		expected := TweetChange{
			ID: "docker",
			Channel: "container",
		}
		actual, err := parseTweetCommand(input)
		assert := assert.New(t)
		assert.Equal(expected, actual)
		assert.Nil(err)
	})

	t.Run("Show command", func(t *testing.T) {
		input := []string{"<@!1234567890>", "tweet", "show"}
		expected := TweetShow{}
		actual, err := parseTweetCommand(input)
		assert := assert.New(t)
		assert.Equal(expected, actual)
		assert.Nil(err)
	})
}
