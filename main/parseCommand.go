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

	switch arguments[1] {
	case "help":
		return Help{}, nil
	case "ping":
		return Ping{}, nil
	case "tweet":
		parseTweetCommand(arguments)
	}
	return nil, errors.New("そのようなコマンドはありません")
}

func parseTweetCommand(arguments []string) (Command, error) {
	if len(arguments) < 3 {
		return nil, errors.New("サブコマンドを指定してください")
	}

	switch arguments[2] {
	case "create":
		if err := assertArguments(arguments, 5, "IDとキーワードを指定してください"); err != nil {
			return nil, err
		}
		return TweetCreate {
			ID: arguments[3],
			Keywords: arguments[4:],
		}, nil
	case "add":
		if err := assertArguments(arguments, 5, "IDとキーワードを指定してください"); err != nil {
			return nil, err
		}
		return TweetAdd {
			ID: arguments[3],
			Keywords: arguments[4:],
		}, nil
	case "remove":
		if err := assertArguments(arguments, 5, "IDとキーワードを指定してください"); err != nil {
			return nil, err
		}
		return TweetRemove {
			ID: arguments[3],
			Keywords: arguments[4:],
		}, nil
	case "delete":
		if err := assertArguments(arguments, 4, "IDを指定してください"); err != nil {
			return nil, err
		}
		return TweetDelete {
			ID: arguments[3],
		}, nil
	case "change":
		if err := assertArguments(arguments, 5, "IDとキーワードを指定してください"); err != nil {
			return nil, err
		}
		channel := strings.Trim(arguments[4], "#")
		return TweetChange {
			ID: arguments[3],
			Channel: channel,
		}, nil
	case "show":
		return TweetShow{}, nil
	default:
		return nil, errors.New("適切なサブコマンドを指定してください")
	}
}

// assertArguments tests whether the number of elements in a `arguments` is more than `length` or not.
// If not, return error originated with `message`.
func assertArguments(arguments []string, length int, message string) error {
	if len(arguments) < length {
		return errors.New(message)
	}
	return nil
}
