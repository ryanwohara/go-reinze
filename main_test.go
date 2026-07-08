package main

import (
	"testing"

	irc "github.com/thoj/go-ircevent"
)

func Test_SafelyRecoversPanic(t *testing.T) {
	safely("test", func() {
		panic("boom")
	})
}

func Test_HandlePrivmsgEmptyMessage(t *testing.T) {
	// PRIVMSG #chan : from Dragon must not panic on Message()[1:]
	handlePrivmsg(nil, &irc.Event{Nick: "Dragon", Arguments: []string{"#chan", ""}})
}

func Test_HandlePrivmsgNoArguments(t *testing.T) {
	handlePrivmsg(nil, &irc.Event{Nick: "Dragon", Arguments: []string{}})
}

func Test_HandlePrivmsgUnknownCommand(t *testing.T) {
	handlePrivmsg(nil, &irc.Event{Nick: "Dragon", Arguments: []string{"#chan", "+notacommand"}})
}

func Test_HandlePrivmsgOtherNick(t *testing.T) {
	handlePrivmsg(nil, &irc.Event{Nick: "someone", Arguments: []string{"#chan", "+players"}})
}

func Test_HandleInviteMissingArguments(t *testing.T) {
	handleInvite(nil, &irc.Event{Nick: "Dragon", Arguments: []string{"nick"}})
}

func Test_HandleJoinMissingArguments(t *testing.T) {
	handleJoin(nil, &irc.Event{Nick: "someone", Arguments: []string{}})
}

func Test_HandleJoinOtherChannel(t *testing.T) {
	handleJoin(nil, &irc.Event{Nick: "someone", Arguments: []string{"#not-the-greet-channel"}})
}
