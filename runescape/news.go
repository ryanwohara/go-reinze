package runescape

import (
	"database/sql"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"strings"
	irc "github.com/thoj/go-ircevent"
)

func checkNews(db *sql.DB, irccon *irc.Connection) {
	rs3 := acquireRs3News(db)
	rs3 = append(rs3, generateHash(rs3), "runescape3")
	rs3Exists := queryExists(db, rs3)

	osrs := acquireOsrsNews(db)
	osrs = append(osrs, generateHash(osrs), "oldschool")
	osrsExists := queryExists(db, osrs)

	fmt.Println(rs3Exists, osrsExists)

	if osrsExists == 0 {
		writeNewsToDb(osrs)
	}
	if rs3Exists == 0 {
		writeNewsToDb(rs3)
	}
	if osrsExists == 0 || rs3Exists == 0 {
		updateTopic(constructTopic(rs3, osrs), irccon)
	}
}

func queryExists(db *sql.DB, rs []string) int {
	var exists int

	db.QueryRow("SELECT exists (SELECT hash_id FROM `rsnews` WHERE hash_id = ?)", rs[2]).Scan(&exists)

	return exists
}

func constructTopic(rs3 []string, osrs []string) string {
	message := fmt.Sprintf("3https://rshelp.com | https://discord.gg/xAVtSdMdhU | RS3: 04%s:04 06https://rshelp.com/t/%s 03| OSRS: 04%s:04 06https://rshelp.com/t/%s", rs3[0], rs3[2], osrs[0], osrs[2])

	return message
}

func updateTopic(topic string, irccon *irc.Connection) {
	// irccon.SendRawf("TOPIC #rshelp :%s", topic)
	irccon.SendRawf("PRIVMSG #asdfghj :%s", topic)
}

func acquireRs3News(db *sql.DB) []string {
	res, err := http.Get("https://www.runescape.com/community")

	maybePanic(err)

	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	maybePanic(err)

	var articles [][]string

	doc.Find("article").Each(func(i int, s *goquery.Selection) {
		title := strings.Replace(s.Find("h4 a").Text(), "This Week in RuneScape", "TWIR", 1)
		link, _ := s.Find("h4 a").Attr("href")

		articles = append(articles, []string{title, link})
	})

	return articles[0]
}

func acquireOsrsNews(db *sql.DB) []string {
	res, err := http.Get("https://oldschool.runescape.com")

	maybePanic(err)

	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	maybePanic(err)

	var articles [][]string

	doc.Find("article").Each(func(i int, s *goquery.Selection) {
		title := strings.Replace(s.Find("h3 a").Text(), "This Week In RuneScape", "TWIR", 1)
		link, _ := s.Find("h3 a").Attr("href")

		articles = append(articles, []string{title, link})
	})

	return articles[0]
}

func generateHash(news []string) string {
	hash := getHash(news[0] + news[1])

	return hash
}

func getHash(url string) string {
	hash := sha256.Sum256([]byte(url))
	return hex.EncodeToString(hash[:])[:5]
}

func writeNewsToDb(news []string) {
	db := db()

	db.Ping()

	_, err := db.Exec("INSERT INTO `rsnews` (title, url, hash_id, runescape) VALUES (?, ?, ?, ?)", news[0], news[1], news[2], news[3])
	maybePanic(err)

	db.Close()
}

type News struct {
	Title     string `db:"title" json:"title"`
	Url       string `db:"url" json:"url"`
	Hash      string `db:"hash_id" json:"hash_id"`
	Runescape string `db:"runescape" json:"runescape"`
}
