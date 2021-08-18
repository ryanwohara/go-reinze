package runescape

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"strings"
)

func checkNews() string {
	rs3 := rs3()
	rs3 = append(rs3, storeNews(rs3[0], rs3[1], "runescape3"))

	osrs := osrs()
	osrs = append(osrs, storeNews(osrs[0], osrs[1], "oldschool"))

	return constructTopic(rs3, osrs)
}

func rs3() []string {
	return acquireRs3News()
}

func osrs() []string {
	return acquireOsrsNews()
}

func constructTopic(rs3 []string, osrs []string) string {
	message := fmt.Sprintf("3https://rshelp.com | https://discord.gg/xAVtSdMdhU | RS3: 04%s:04 06https://rshelp.com/t/%s 03| OSRS: 04%s:04 06https://rshelp.com/t/%s", rs3[0], rs3[2], osrs[0], osrs[2])

	return message
}

func acquireRs3News() []string {
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
		title := strings.Replace(s.Find("h4 a").Text(), "This Week In RuneScape", "TWIR", 1)
		link, _ := s.Find("h4 a").Attr("href")

		articles = append(articles, []string{title, link})
	})

	return articles[0]
}

func acquireOsrsNews() []string {
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

func storeNews(title string, url string, runescape string) string {
	hash := getHash(title + url)
	writeNewsToDb(title, url, hash, runescape)

	return hash
}

func getHash(url string) string {
	hash := sha256.Sum256([]byte(url))
	return hex.EncodeToString(hash[:])[:5]
}

func writeNewsToDb(title string, url string, hash string, runescape string) {
	db := db()

	db.Ping()

	var exists bool

	db.QueryRow("SELECT exists (SELECT hash_id FROM `rsnews` WHERE title = ? AND url = ? AND hash_id = ? AND runescape = ?)", title, url, hash, runescape).Scan(&exists)

	if !exists {
		_, err := db.Exec("INSERT INTO `rsnews` (title, url, hash_id, runescape) VALUES (?, ?, ?, ?)", title, url, hash, runescape)
		maybePanic(err)
	}

	db.Close()
}

type News struct {
	Title     string `db:"title" json:"title"`
	Url       string `db:"url" json:"url"`
	Hash      string `db:"hash_id" json:"hash_id"`
	Runescape string `db:"runescape" json:"runescape"`
}
