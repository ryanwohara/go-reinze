package news

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"

	rss "github.com/mmcdole/gofeed"
	irc "github.com/thoj/go-ircevent"
)

func CheckNews(db *sql.DB, irccon *irc.Connection) {
	sources := [1]string{
		"https://www.denverpost.com/feed/", // https://www.denverpost.com/web-feeds/
	}

	for _, url := range sources {
		feedparser := rss.NewParser()
		feed, err := feedparser.ParseURL(url)

		if err == nil {
			for _, item := range feed.Items {
				news := News{
					Title: item.Title,
					Url:   item.Link,
					Hash:  generateHash(item),
				}
				processNews(db, irccon, feed, news)
			}
		}
	}
}

func processNews(db *sql.DB, irccon *irc.Connection, feed *rss.Feed, news News) {
	if !queryExists(db, news) {
		if writeNewsToDb(db, news) {
			irccon.SendRawf("PRIVMSG #news :[%s] %s [%s]", feed.Title, news.Title, news.Url)
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
	var hash_id string

	err := db.QueryRow("SELECT hash_id FROM `news` WHERE hash_id = ?", news.Hash).Scan(&hash_id)

	return (err == nil)
}

func writeNewsToDb(db *sql.DB, news News) bool {
	_, err := db.Exec("INSERT INTO `news` (title, url, hash_id) VALUES (?, ?, ?)", news.Title, news.Url, news.Hash)

	return (err == nil)
}
