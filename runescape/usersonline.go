package runescape

import (
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"strconv"
)

func getUsersOnline() []int {
	resp, err := http.Get("https://oldschool.runescape.com/")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	re := regexp.MustCompile("<p class='player-count'>There are currently ([0-9,]+) people playing!</p>")
	result := re.FindString(string(body))
	players := strings.Fields(result)[4]
	osrs, _ := strconv.Atoi(strings.Replace(players, ",", "", -1))

	resp, err = http.Get("https://www.runescape.com/player_count.js?varname=iPlayerCount&callback=jQuery36004811633109689837_1628665230298")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	re = regexp.MustCompile("([0-9]+)")
	result2 := re.FindAllString(string(body), 3)
	players = result2[2]
	total, _ := strconv.Atoi(players)
	rs3 := total - osrs

	return []int{osrs, rs3, total}
}
