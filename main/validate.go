package main

import (
	"errors"
	"regexp"
	"strings"
)

const idPatterns = `^[a-zA-Z0-9_]+$`

var idRegexp = regexp.MustCompile(idPatterns)

// ValidateTwitterID checks if given ID is well formed.
func ValidateTwitterID(id string) (string, error) {
	if idRegexp.MatchString(id) {
		return id, nil
	}
	return "", errors.New("TwitterのIDには英数字とアンダーバーのみ使えます")
}

// FormatChannelID trims "#" to normalize channel ID of Discord.
func FormatChannelID(id string) string {
	return strings.Trim(id, "#")
}
