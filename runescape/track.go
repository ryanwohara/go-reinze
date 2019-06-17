package runescape

import (
	"fmt"
)

// Track will check CrystalMathLabs
// for updated tracking data for a
// given RuneScape name.
func Track(message string) string {
	return track(message)
}

func track(message string) string {
	fmt.Println(message)
	return "no data for you"
}
