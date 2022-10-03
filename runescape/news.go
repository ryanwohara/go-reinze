package runescape

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	irc "github.com/thoj/go-ircevent"
)

func checkNews(db *sql.DB, irccon *irc.Connection) {
	rs3 := getNews("https://www.runescape.com/community", "h4 a")
	rs3Exists := queryExists(db, rs3)
	if !rs3Exists {
		writeNewsToDb(db, rs3)
	}

	osrs := getNews("https://oldschool.runescape.com", "h3 a")
	osrsExists := queryExists(db, osrs)
	if !osrsExists {
		writeNewsToDb(db, osrs)
	}

	if !osrsExists || !rs3Exists {
		updateTopic(constructTopic(rs3, osrs), irccon)
	}
}

func queryExists(db *sql.DB, rs []string) bool {
	var hash_id string

	err := db.QueryRow("SELECT hash_id FROM `rsnews` WHERE hash_id = ?", rs[2]).Scan(&hash_id)

	return (err == nil)
}

func constructTopic(rs3 []string, osrs []string) string {
	message := fmt.Sprintf("3https://rshelp.com | https://discord.gg/xAVtSdMdhU | RS3: 04%s:04 06https://rshelp.com/t/%s 03| OSRS: 04%s:04 06https://rshelp.com/t/%s", rs3[0], rs3[2], osrs[0], osrs[2])

	return message
}

func updateTopic(topic string, irccon *irc.Connection) {
	irccon.SendRawf("TOPIC #rshelp :%s", topic)
}

func getNews(url string, element string) []string {
	res, err := http.Get(url)

	if err != nil {
		return []string{}
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return []string{}
	}

	var article []string

	doc.Find("article").First().Each(func(i int, s *goquery.Selection) {
		title := strings.Replace(s.Find(element).Text(), "This Week In RuneScape", "TWIR", 1)
		link, _ := s.Find(element).Attr("href")

		article = []string{title, link}
	})

	version := getVersion(article[1])
	article = append(article, generateHash(article), version)

	return article
}

func getVersion(url string) string {
	var version string

	if strings.Contains(url, "oldschool") {
		version = "oldschool"
	} else {
		version = "runescape3"
	}

	return version
}

func generateHash(news []string) string {
	hash := getHash(news[0] + news[1])

	return hash
}

func getHash(url string) string {
	hash := sha256.Sum256([]byte(url))

	return hex.EncodeToString(hash[:])[:5]
}

func writeNewsToDb(db *sql.DB, news []string) {
	db.Exec("INSERT INTO `rsnews` (title, url, hash_id, runescape) VALUES (?, ?, ?, ?)", news[0], news[1], news[2], news[3])
}

type News struct {
	Title     string `db:"title" json:"title"`
	Url       string `db:"url" json:"url"`
	Hash      string `db:"hash_id" json:"hash_id"`
	Runescape string `db:"runescape" json:"runescape"`
}

func handleException(err error) {
	fmt.Println(err)
}
