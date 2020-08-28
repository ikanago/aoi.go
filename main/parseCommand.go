package main

import (
	"errors"
	"strings"
)

// ParseCommand parses messages from Discord and returns results as sturct.
// Assumes first word of the input as a mention to this bot.
func ParseCommand(input string) (Command, error) {
	arguments := strings.Fields(input)
	if len(arguments) < 2 {
		return nil, errors.New("コマンドを指定してください")
	}

	if arguments[1] == "help" {
		return Help{}, nil
	} else if arguments[1] == "ping" {
		return Ping{}, nil
	}
	return nil, errors.New("そのようなコマンドはありません")
}
