package runescape

import (
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

func getUsersOnline(getterFunc PageGetter) string {
	url := "https://oldschool.runescape.com/"

	body := string(getterFunc(url))

	re := regexp.MustCompile("<p class='player-count'>There are currently ([0-9,]+) people playing!</p>")
	result := re.FindString(string(body))

	players := strings.Fields(result)

	if len(players) == 0 {
		return ""
	}

	return players[4]
}

// PageGetter is a type to allow mocking
// of the calls to Jagex's website.
type PageGetter func(url string) []byte

func getOsRsPlayerData(url string) []byte {
	resp, err := http.Get(url)
	if err != nil {
		return []byte("")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte("")
	}

	return body
}
