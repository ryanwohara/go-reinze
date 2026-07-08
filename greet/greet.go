package greet

import (
	"math/rand"
	"os"
	"strings"

	irc "github.com/thoj/go-ircevent"
)

func Greet(irccon *irc.Connection, channel string, nick string) {
	greet := pickGreet(os.Getenv("GREET_MESSAGES"), nick)

	if len(greet) == 0 {
		return
	}

	irccon.Action(channel, greet)
}

func pickGreet(messages string, nick string) string {
	var greets []string

	for _, greet := range strings.Split(messages, "\n") {
		if strings.TrimSpace(greet) != "" {
			greets = append(greets, greet)
		}
	}

	if len(greets) == 0 {
		return ""
	}

	return strings.Replace(greets[rand.Intn(len(greets))], "!nick!", nick, -1)
}
