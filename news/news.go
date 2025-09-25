package news

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	rss "github.com/mmcdole/gofeed"
)

type Config struct {
	Target  string   `json:"target"`
	Sources []string `json:"sources"`
}

func CheckNews(db *sql.DB, queue chan string) {
	newsConfig := os.Getenv("NEWS_CONFIG")

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
			feedparser := rss.NewParser()
			feed, err := feedparser.ParseURL(url)

			if err != nil {
				fmt.Println("news.go: " + err.Error() + " " + url)
				continue
			}

			for _, item := range feed.Items {
				if len(item.Link) == 0 {
					item.Link = item.GUID
				}

				news := News{
					Title: item.Title,
					Url:   item.Link,
					Hash:  generateHash(item),
				}

				go processNews(db, target, feed, news, queue)
			}
		}
	}
}

func processNews(db *sql.DB, target string, feed *rss.Feed, news News, queue chan string) {
	if writeNewsToDb(db, news) {
		queue <- fmt.Sprintf("PRIVMSG %s :[%s] %s [%s]", target, feed.Title, news.Title, news.Url)
	}
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
	var count string

	err := db.QueryRow("SELECT COUNT(url) FROM `news` WHERE url = '?'", news.Url).Scan(&count)
	if err != nil {
		fmt.Println("news/news.go: queryExists QueryRow:" + err.Error())
		return true // we'll return true to prevent messages being sent to the network
	}

	countInt, err := strconv.Atoi(count)
	if err != nil {
		fmt.Println("news/news.go: queryExists strconv:" + err.Error())
		return true
	}

	return countInt > 0
}

func writeNewsToDb(db *sql.DB, news News) bool {
	if queryExists(db, news) {
		return false
	}

	_, err := db.Exec("INSERT INTO `news` (title, url, hash_id) VALUES ('?', '?', '?')", news.Title, news.Url, news.Hash)

	if err != nil {
		fmt.Println("news/news.go: writeNewsToDb:" + err.Error())
	}

	return err == nil
}
