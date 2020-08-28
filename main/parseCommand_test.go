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
