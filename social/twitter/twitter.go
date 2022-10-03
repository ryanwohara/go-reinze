package twitter

import (
	"context"
	"database/sql"
	"os"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	irc "github.com/thoj/go-ircevent"
	"golang.org/x/oauth2/clientcredentials"
)

var (
	ConsumerKey    string = os.Getenv("TWITTER_CONSUMER_KEY")
	ConsumerSecret string = os.Getenv("TWITTER_CONSUMER_SECRET")
)

func CheckPosts(db *sql.DB, irccon *irc.Connection) {
	checkPosts(db, irccon, "stonrus", "#stonr")
}

func checkPosts(db *sql.DB, irccon *irc.Connection, source string, destination string) {
	config := &clientcredentials.Config{
		ClientID:     ConsumerKey,
		ClientSecret: ConsumerSecret,
		TokenURL:     "https://api.twitter.com/oauth2/token",
	}
	httpClient := config.Client(context.Background())

	client := twitter.NewClient(httpClient)

	var params twitter.UserTimelineParams
	params.ScreenName = source

	tweets, _, err := client.Timelines.UserTimeline(&params)

	if err == nil {
		for _, tweet := range tweets {
			if !queryExists(db, tweet) {
				writeTweetToDb(db, tweet)

				irccon.SendRawf("PRIVMSG %s :[%s Twitter] %s", destination, source, tweet.Text)

				time.Sleep(2 * time.Second)
			}
		}
	}
}

func queryExists(db *sql.DB, tweet twitter.Tweet) bool {
	var tweet_id string

	err := db.QueryRow("SELECT tweet_id FROM `twitter` WHERE tweet_id = ?", tweet.ID).Scan(&tweet_id)

	return (err == nil)
}

func writeTweetToDb(db *sql.DB, tweet twitter.Tweet) bool {
	_, err := db.Exec("INSERT INTO `twitter` (text, tweet_id) VALUES (?, ?)", tweet.Text, tweet.ID)

	return (err == nil)
}
