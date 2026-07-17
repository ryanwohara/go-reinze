package news

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	rss "github.com/mmcdole/gofeed"
)

type Config struct {
	Target  string   `json:"target"`
	Sources []string `json:"sources"`
}

func CheckNews(db *sql.DB, queue chan string) {
	newsConfig := os.Getenv("NEWS_CONFIG")
	feedparser := rss.NewParser()
	feedparser.UserAgent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/140.0.0.0 Safari/537.36"
	feedparser.Client = &http.Client{Timeout: 30 * time.Second}

	var structConfig []Config

	err := json.Unmarshal([]byte(newsConfig), &structConfig)

	if err != nil {
		fmt.Println("news.go: " + err.Error())
		return
	}

	for _, config := range structConfig {
		target := config.Target
		sources := config.Sources

		for _, url := range sources {
			feed, err := feedparser.ParseURL(url)

			if err != nil {
				fmt.Println("news.go: " + err.Error() + " " + url)
				continue
			}

			for _, item := range feed.Items {
				processNews(db, target, feed, newsFromItem(item), queue)
			}
		}
	}
}

// The news table caps title and url at VARCHAR(255); MySQL strict mode
// rejects longer values outright, so clamp them here — before the dedup
// query — so the stored row and the existence check always agree.
const maxColumnLength = 255

func newsFromItem(item *rss.Item) News {
	if len(item.Link) == 0 {
		item.Link = item.GUID
	}

	return News{
		Title: truncateRunes(item.Title, maxColumnLength),
		Url:   truncateRunes(stripGetParams(item.Link), maxColumnLength),
		Hash:  generateHash(item),
	}
}

func truncateRunes(s string, max int) string {
	runes := []rune(s)
	if len(runes) <= max {
		return s
	}

	return string(runes[:max])
}

func processNews(db *sql.DB, target string, feed *rss.Feed, news News, queue chan string) {
	if writeNewsToDb(db, news) {
		queue <- fmt.Sprintf("PRIVMSG %s :[%s] %s [%s]", target, feed.Title, news.Title, news.Url)
	}
}

func stripGetParams(url string) string {
	return strings.Split(url, "?")[0]
}

func generateHash(item *rss.Item) string {
	return getHash(item.Title + "~" + item.Link)
}

func getHash(toHash string) string {
	hash := sha256.Sum256([]byte(toHash))

	return hex.EncodeToString(hash[:])[:5]
}

type News struct {
	Title string `db:"title" json:"title"`
	Url   string `db:"url" json:"url"`
	Hash  string `db:"hash_id" json:"hash_id"`
}

func queryExists(db *sql.DB, news News) bool {
	var count int

	err := db.QueryRow("SELECT COUNT(url) FROM `news` WHERE url = ?", news.Url).Scan(&count)
	if err != nil {
		fmt.Println("news/news.go: queryExists QueryRow:" + err.Error())
		return true // we'll return true to prevent messages being sent to the network
	}

	return count > 0
}

func writeNewsToDb(db *sql.DB, news News) bool {
	if queryExists(db, news) {
		return false
	}

	_, err := db.Exec("INSERT INTO `news` (title, url, hash_id) VALUES (?, ?, ?)", news.Title, news.Url, news.Hash)

	if err != nil {
		fmt.Println("news/news.go: writeNewsToDb:" + err.Error())
	}

	return err == nil
}
