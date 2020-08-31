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
		return parseTweetCommand(arguments)
	case "memo":
		return parseMemoCommand(arguments)
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
		screenName, err := ValidateTwitterID(arguments[3])
		if err != nil {
			return nil, err
		}
		return TweetCreate{
			ScreenName: screenName,
			Keywords:   arguments[4:],
		}, nil
	case "add":
		if err := assertArguments(arguments, 5, "IDとキーワードを指定してください"); err != nil {
			return nil, err
		}
		screenName, err := ValidateTwitterID(arguments[3])
		if err != nil {
			return nil, err
		}
		return TweetAdd{
			ScreenName: screenName,
			Keywords:   arguments[4:],
		}, nil
	case "remove":
		if err := assertArguments(arguments, 5, "IDとキーワードを指定してください"); err != nil {
			return nil, err
		}
		screenName, err := ValidateTwitterID(arguments[3])
		if err != nil {
			return nil, err
		}
		return TweetRemove{
			ScreenName: screenName,
			Keywords:   arguments[4:],
		}, nil
	case "delete":
		if err := assertArguments(arguments, 4, "IDを指定してください"); err != nil {
			return nil, err
		}
		screenName, err := ValidateTwitterID(arguments[3])
		if err != nil {
			return nil, err
		}
		return TweetDelete{
			ScreenName: screenName,
		}, nil
	case "change":
		if err := assertArguments(arguments, 5, "IDとキーワードを指定してください"); err != nil {
			return nil, err
		}
		screenName, err := ValidateTwitterID(arguments[3])
		if err != nil {
			return nil, err
		}
		return TweetChange{
			ScreenName: screenName,
			Channel:    FormatChannelID(arguments[4]),
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

func parseMemoCommand(arguments []string) (Command, error) {
	if len(arguments) < 3 {
		return nil, errors.New("メモしたい発言またはサブコマンドを指定してください")
	}

	if arguments[2] == "show" {
		return MemoShow{}, nil
	}
	memo := strings.Join(arguments[2:], " ")
	return MemoRegister{
		Text: memo,
	}, nil
}
