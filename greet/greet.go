package greet

import (
	"math/rand"
	"os"
	"strings"

	irc "github.com/thoj/go-ircevent"
)

func Greet(irccon *irc.Connection, channel string, nick string) {
	greets := strings.Split(os.Getenv("GREET_MESSAGES"), "\n")

	greet := strings.Replace(greets[rand.Intn(len(greets)-1)], "!nick!", nick, -1)

	irccon.Action(channel, greet)
}
