package runescape

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var httpClient = &http.Client{Timeout: 30 * time.Second}

type TotalPlayers struct {
	Accounts          int    `json:"accounts"`
	AccountsFormatted string `json:"accountsformatted"`
}

func getHttpContent(url string) (string, error) {
	resp, err := httpClient.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Printf("status code error: %d %s", resp.StatusCode, resp.Status)
		return "", errors.New(strconv.Itoa(resp.StatusCode))
	}
	body, err := io.ReadAll(resp.Body)
	return string(body), err

}

func getUsersOnline() []int {
	return getUsersOnlineFrom(
		"https://oldschool.runescape.com/",
		"https://www.runescape.com/player_count.js?varname=iPlayerCount&callback=jQuery36004811633109689837_1628665230298",
		"https://secure.runescape.com/m=account-creation-reports/rsusertotal.ws",
	)
}

func getUsersOnlineFrom(osrsUrl string, playerCountUrl string, userTotalUrl string) []int {
	body, err := getHttpContent(osrsUrl)

	if err != nil {
		return []int{0, 0, 0, 0}
	}

	re := regexp.MustCompile("<p class='player-count'>There are currently ([0-9,]+) people playing!</p>")
	result := re.FindString(string(body))
	fields := strings.Fields(result)
	if len(fields) < 5 {
		return []int{0, 0, 0, 0}
	}
	players := fields[4]
	osrs, _ := strconv.Atoi(strings.Replace(players, ",", "", -1))

	body, err = getHttpContent(playerCountUrl)

	if err != nil {
		return []int{0, 0, 0, 0}
	}

	re = regexp.MustCompile("([0-9]+)")
	result2 := re.FindAllString(string(body), 3)
	if len(result2) < 3 {
		return []int{0, 0, 0, 0}
	}
	players = result2[2]
	total_online, _ := strconv.Atoi(players)
	rs3 := total_online - osrs

	body, err = getHttpContent(userTotalUrl)

	if err != nil {
		return []int{0, 0, 0, 0}
	}

	var t TotalPlayers
	err = json.Unmarshal([]byte(body), &t)

	if err != nil {
		return []int{0, 0, 0, 0}
	}

	total_accounts := t.Accounts

	return []int{osrs, rs3, total_online, total_accounts}
}
