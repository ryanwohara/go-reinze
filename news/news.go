package news

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"os"
	"time"

	rss "github.com/mmcdole/gofeed"
	irc "github.com/thoj/go-ircevent"
)

func CheckNews(db *sql.DB, irccon *irc.Connection) {
	sources := []string{
		"https://www.denverpost.com/feed/",
		"http://rss.slashdot.org/Slashdot/slashdotMain",
		"https://www.cbc.ca/cmlink/rss-topstories",
		"https://www.techradar.com/rss",
		"https://news.yahoo.com/rss",
		"https://www.majorgeeks.com/files/rss",
		"http://rss.cnn.com/rss/cnn_topstories.rss",
		"https://moxie.foxnews.com/google-publisher/latest.xml",
		"https://wsvn.com/feed/",
		"https://wgntv.com/feed/",
		"https://theatlantavoice.com/feed/",
		"https://nerdist.com/feed/",
		"https://feeds.skynews.com/feeds/rss/world.xml",
		"https://www.latimes.com/news/rss2.0.xml",
		"http://feeds.bbci.co.uk/news/world/rss.xml",
	}

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

				processNews(db, irccon, feed, news)
			}
		}
	}
}

func processNews(db *sql.DB, irccon *irc.Connection, feed *rss.Feed, news News) {
	if !queryExists(db, news) {
		if writeNewsToDb(db, news) {
			irccon.SendRawf("PRIVMSG %s :[%s] %s [%s]", os.Getenv("NEWS_CHANNEL"), feed.Title, news.Title, news.Url)
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
	var hash_id string

	err := db.QueryRow("SELECT hash_id FROM `news` WHERE hash_id = ?", news.Hash).Scan(&hash_id)

	return (err == nil)
}

func writeNewsToDb(db *sql.DB, news News) bool {
	_, err := db.Exec("INSERT INTO `news` (title, url, hash_id) VALUES (?, ?, ?)", news.Title, news.Url, news.Hash)

	return (err == nil)
}
