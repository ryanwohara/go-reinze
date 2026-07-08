package greet

import (
	"testing"
)

func Test_PickGreetEmptyConfig(t *testing.T) {
	greet := pickGreet("", "Bob")

	if greet != "" {
		t.Errorf("Expecting empty greeting, received %q", greet)
	}
}

func Test_PickGreetSingleMessage(t *testing.T) {
	greet := pickGreet("hello !nick!", "Bob")

	if greet != "hello Bob" {
		t.Errorf("Expecting %q, received %q", "hello Bob", greet)
	}
}

func Test_PickGreetSelectsEveryMessage(t *testing.T) {
	seen := map[string]bool{}

	for i := 0; i < 200; i++ {
		seen[pickGreet("hi !nick!\nyo !nick!", "Bob")] = true
	}

	if !seen["hi Bob"] || !seen["yo Bob"] {
		t.Errorf("Expecting both greetings to be selectable, received %v", seen)
	}
}

func Test_PickGreetTrailingNewline(t *testing.T) {
	for i := 0; i < 50; i++ {
		greet := pickGreet("hi !nick!\n", "Bob")

		if greet != "hi Bob" {
			t.Errorf("Expecting %q, received %q", "hi Bob", greet)
		}
	}
}
