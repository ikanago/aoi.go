package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRemoveElement(t *testing.T) {
	t.Run("Empty array", func(t *testing.T) {
		actual := removeElement([]string{}, []string{})
		assert := assert.New(t)
		assert.Equal([]string(nil), actual)
	})

	t.Run("No entry", func(t *testing.T) {
		expected := []string{"hoge", "fuga"}
		actual := removeElement([]string{"hoge", "fuga"}, []string{"neko"})
		assert := assert.New(t)
		assert.Equal(expected, actual)
	})

	t.Run("Delete an element", func(t *testing.T) {
		expected := []string{"hoge"}
		actual := removeElement([]string{"hoge", "fuga"}, []string{"fuga"})
		assert := assert.New(t)
		assert.Equal(expected, actual)
	})

	t.Run("Delete elements", func(t *testing.T) {
		expected := []string{"hoge", "neko"}
		actual := removeElement(
			[]string{"hoge", "fuga", "neko", "bar"},
			[]string{"fuga", "bar"},
		)
		assert := assert.New(t)
		assert.Equal(expected, actual)
	})
}
