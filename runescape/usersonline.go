package runescape

import (
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"strconv"
	"encoding/json"
)

type TotalPlayers struct {
	Accounts int `json:"accounts"`
	AccountsFormatted string `json:"accountsformatted"`
}

func maybePanic(err error) {
	if err != nil {
		panic(err)
	}
}

func getHttpContent(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return string(body), err
}

func getUsersOnline() []int {
	body, err := getHttpContent("https://oldschool.runescape.com/")
	maybePanic(err)
	re := regexp.MustCompile("<p class='player-count'>There are currently ([0-9,]+) people playing!</p>")
	result := re.FindString(string(body))
	players := strings.Fields(result)[4]
	osrs, _ := strconv.Atoi(strings.Replace(players, ",", "", -1))

	body, err = getHttpContent("https://www.runescape.com/player_count.js?varname=iPlayerCount&callback=jQuery36004811633109689837_1628665230298")
	maybePanic(err)
	re = regexp.MustCompile("([0-9]+)")
	result2 := re.FindAllString(string(body), 3)
	players = result2[2]
	total_online, _ := strconv.Atoi(players)
	rs3 := total_online - osrs

	body, err = getHttpContent("https://secure.runescape.com/m=account-creation-reports/rsusertotal.ws")
	maybePanic(err)
	var t TotalPlayers
	err = json.Unmarshal([]byte(body), &t)
	maybePanic(err)
	total_accounts := t.Accounts

	return []int{osrs, rs3, total_online, total_accounts}
}
