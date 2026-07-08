package runescape

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

const deadUrl = "http://127.0.0.1:1/"

func newsServer(body string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, body)
	}))
}

func Test_GetUsersOnline(t *testing.T) {
	osrs := newsServer("<html><body><p class='player-count'>There are currently 123,456 people playing!</p></body></html>")
	defer osrs.Close()

	playerCount := newsServer("jQuery36004811633109689837_1628665230298(210000);")
	defer playerCount.Close()

	userTotal := newsServer(`{"accounts":300000000,"accountsformatted":"300,000,000"}`)
	defer userTotal.Close()

	players := getUsersOnlineFrom(osrs.URL, playerCount.URL, userTotal.URL)

	expected := []int{123456, 86544, 210000, 300000000}
	for i := range expected {
		if players[i] != expected[i] {
			t.Errorf("Expecting %v, received %v", expected, players)
			break
		}
	}
}

func Test_GetUsersOnlineUnreachableHost(t *testing.T) {
	players := getUsersOnlineFrom(deadUrl, deadUrl, deadUrl)

	assertZeros(t, players)
}

func Test_GetUsersOnlineNon200(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer srv.Close()

	players := getUsersOnlineFrom(srv.URL, srv.URL, srv.URL)

	assertZeros(t, players)
}

func Test_GetUsersOnlineNoPlayerCountMarkup(t *testing.T) {
	osrs := newsServer("<html><body><p>maintenance</p></body></html>")
	defer osrs.Close()

	players := getUsersOnlineFrom(osrs.URL, deadUrl, deadUrl)

	assertZeros(t, players)
}

func Test_GetUsersOnlineMalformedPlayerCountJs(t *testing.T) {
	osrs := newsServer("<html><body><p class='player-count'>There are currently 123,456 people playing!</p></body></html>")
	defer osrs.Close()

	playerCount := newsServer("jQuery(1);")
	defer playerCount.Close()

	players := getUsersOnlineFrom(osrs.URL, playerCount.URL, deadUrl)

	assertZeros(t, players)
}

func Test_GetUsersOnlineBadUserTotalJson(t *testing.T) {
	osrs := newsServer("<html><body><p class='player-count'>There are currently 123,456 people playing!</p></body></html>")
	defer osrs.Close()

	playerCount := newsServer("jQuery36004811633109689837_1628665230298(210000);")
	defer playerCount.Close()

	userTotal := newsServer("not json")
	defer userTotal.Close()

	players := getUsersOnlineFrom(osrs.URL, playerCount.URL, userTotal.URL)

	assertZeros(t, players)
}

func Test_FormatUsersOnline(t *testing.T) {
	message := formatUsersOnline([]int{100, 100, 200, 1000})

	expected := "There are currently 100 OSRS players (50.00%) and 100 RS3 players (50.00%) online. (Total: 200) (Total Registered Accounts: 1,000)"
	if message != expected {
		t.Errorf("Expecting %q, received %q", expected, message)
	}
}

func Test_FormatUsersOnlineOutage(t *testing.T) {
	message := formatUsersOnline([]int{0, 0, 0, 0})

	expected := "Player counts are currently unavailable."
	if message != expected {
		t.Errorf("Expecting %q, received %q", expected, message)
	}
}

func assertZeros(t *testing.T, players []int) {
	t.Helper()

	if len(players) != 4 {
		t.Fatalf("Expecting 4 values, received %v", players)
	}

	for i, p := range players {
		if p != 0 {
			t.Errorf("Expecting all zeros, received %v at index %d", players, i)
		}
	}
}
