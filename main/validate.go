package main

import (
	"errors"
	"regexp"
	"strings"
)

const idPatterns = `^[a-zA-Z0-9_]+$`

var idRegexp = regexp.MustCompile(idPatterns)

func validateTwitterID(id string) (string, error) {
	if idRegexp.MatchString(id) {
		return id, nil
	}
	return "", errors.New("TwitterのIDには英数字とアンダーバーのみ使えます")
}

func formatChannleID(id string) string {
	return strings.Trim(id, "#")
}
