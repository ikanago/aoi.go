package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidateTwitterID(t *testing.T) {
	t.Run("Valid ID", func(t *testing.T) {
		input := "ab_12_cd"
		expected := "ab_12_cd"
		actual, err := validateTwitterID(input)
		assert := assert.New(t)
		assert.Equal(expected, actual)
		assert.Nil(err)
	})

	t.Run("Invalid ID", func(t *testing.T) {
		input := "ab_12_cd?"
		expected := ""
		actual, err := validateTwitterID(input)
		assert := assert.New(t)
		assert.Equal(expected, actual)
		assert.NotNil(err)
	})

	t.Run("Invalid ID", func(t *testing.T) {
		input := "あいうえお"
		expected := ""
		actual, err := validateTwitterID(input)
		assert := assert.New(t)
		assert.Equal(expected, actual)
		assert.NotNil(err)
	})
}
