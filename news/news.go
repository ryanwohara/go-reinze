package news

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	rss "github.com/mmcdole/gofeed"
	irc "github.com/thoj/go-ircevent"
)

type Config struct {
	Target  string   `json:"target"`
	Sources []string `json:"sources"`
}

func CheckNews(db *sql.DB, irccon *irc.Connection) {
	news_config := os.Getenv("NEWS_CONFIG")

	var struct_config []Config

	err := json.Unmarshal([]byte(news_config), &struct_config)

	if err == nil {

		for _, config := range struct_config {
			target := config.Target
			sources := config.Sources

			for _, url := range sources {
				feedparser := rss.NewParser()
				feed, err := feedparser.ParseURL(url)

				if err == nil {
					for _, item := range feed.Items {
						if len(item.Link) == 0 {
							item.Link = item.GUID
						}

						news := News{
							Title: item.Title,
							Url:   item.Link,
							Hash:  generateHash(item),
						}

						processNews(db, irccon, target, feed, news)
					}
				}
			}
		}
	} else {
		fmt.Println("news.go: " + err.Error())
	}
}

func processNews(db *sql.DB, irccon *irc.Connection, target string, feed *rss.Feed, news News) {
	if !queryExists(db, news) {
		if writeNewsToDb(db, news) {
			irccon.SendRawf("PRIVMSG %s :[%s] %s [%s]", target, feed.Title, news.Title, news.Url)
			time.Sleep(2 * time.Second)
		}
	}
}

func generateHash(item *rss.Item) string {
	hash := getHash(item.Title + "~" + item.Link)

	return hash
}

func getHash(to_hash string) string {
	hash := sha256.Sum256([]byte(to_hash))

	return hex.EncodeToString(hash[:])[:5]
}

type News struct {
	Title string `db:"title" json:"title"`
	Url   string `db:"url" json:"url"`
	Hash  string `db:"hash_id" json:"hash_id"`
}

func queryExists(db *sql.DB, news News) bool {
	var count string

	err := db.QueryRow("SELECT COUNT(url) FROM `news` WHERE url = ?", news.Url).Scan(&count)

	if err != nil {
		fmt.Println("news/news.go: " + err.Error())
		return true // we'll return true to prevent messages being sent to the network
	}

	count_int, err := strconv.Atoi(count)

	if err != nil {
		fmt.Println("news/news.go: " + err.Error())
		return true
	}

	return (count_int > 0)
}

func writeNewsToDb(db *sql.DB, news News) bool {
	_, err := db.Exec("INSERT INTO `news` (title, url, hash_id) VALUES (?, ?, ?)", news.Title, news.Url, news.Hash)

	return (err == nil)
}
