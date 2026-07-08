package runescape

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
	irc "github.com/thoj/go-ircevent"
)

func checkNews(irccon *irc.Connection, database *sql.DB) {
	rs3 := getNews("https://www.runescape.com/community", "h4 a")
	osrs := getNews("https://oldschool.runescape.com", "h3 a")

	// getNews returns an empty slice when the fetch or scrape fails;
	// skip this run and let the next cron tick retry.
	if len(rs3) < 4 || len(osrs) < 4 {
		log.Printf("news fetch failed (rs3 fields: %d, osrs fields: %d), skipping", len(rs3), len(osrs))
		return
	}

	rs3Exists := queryExists(database, rs3)
	if !rs3Exists {
		writeNewsToDb(database, rs3)
	}

	osrsExists := queryExists(database, osrs)
	if !osrsExists {
		writeNewsToDb(database, osrs)
	}

	if !osrsExists || !rs3Exists {
		updateTopic(constructTopic(rs3, osrs), irccon)
	}
}

func queryExists(db *sql.DB, rs []string) bool {
	var hash_id string

	if len(rs) < 3 {
		return false
	}

	err := db.QueryRow("SELECT hash_id FROM `rsnews` WHERE hash_id = ?", rs[2]).Scan(&hash_id)

	return (err == nil)
}

func constructTopic(rs3 []string, osrs []string) string {
	message := fmt.Sprintf("3https://rshelp.com | https://discord.gg/swiftirc | RS3: 04%s:04 06https://rshelp.com/t/%s 03| OSRS: 04%s:04 06https://rshelp.com/t/%s", rs3[0], rs3[2], osrs[0], osrs[2])

	return message
}

func updateTopic(topic string, irccon *irc.Connection) {
	irccon.SendRawf("TOPIC #rshelp :%s", topic)
}

func getNews(url string, element string) []string {
	res, err := httpClient.Get(url)

	if err != nil {
		return []string{}
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Printf("status code error: %d %s", res.StatusCode, res.Status)
		return []string{}
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

	if len(article) < 2 {
		return []string{}
	}

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
	if len(news) < 4 {
		return
	}

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
