package twitter

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	irc "github.com/thoj/go-ircevent"
	"golang.org/x/oauth2/clientcredentials"
)

var (
	ConsumerKey    string = os.Getenv("TWITTER_CONSUMER_KEY")
	ConsumerSecret string = os.Getenv("TWITTER_CONSUMER_SECRET")
)

type Config struct {
	Source string `json:"source"`
	Target string `json:"target"`
}

func CheckPosts(db *sql.DB, irccon *irc.Connection) {
	twitter_config := os.Getenv("TWITTER_CONFIG")

	var struct_config []Config

	err := json.Unmarshal([]byte(twitter_config), &struct_config)

	if err == nil {
		config := &clientcredentials.Config{
			ClientID:     ConsumerKey,
			ClientSecret: ConsumerSecret,
			TokenURL:     "https://api.twitter.com/oauth2/token",
		}
		httpClient := config.Client(context.Background())

		client := twitter.NewClient(httpClient)

		for _, config := range struct_config {
			checkPosts(db, irccon, client, config.Source, config.Target)
		}
	} else {
		fmt.Println("twitter.go err: " + err.Error())
	}
}

func checkPosts(db *sql.DB, irccon *irc.Connection, client *twitter.Client, source string, destination string) {

	var params twitter.UserTimelineParams
	params.ScreenName = source

	tweets, _, err := client.Timelines.UserTimeline(&params)

	if err == nil {
		for _, tweet := range tweets {
			if !queryExists(db, tweet) {
				writeTweetToDb(db, tweet)

				for _, text := range strings.Split(tweet.Text, "\n") {
					if len(text) > 0 {
						irccon.Privmsgf(destination, "[%s Twitter] %s", source, text)
					}
				}

				time.Sleep(2 * time.Second)
			}
		}
	}
}

func queryExists(db *sql.DB, tweet twitter.Tweet) bool {
	var tweet_id string

	err := db.QueryRow("SELECT tweet_id FROM `twitter` WHERE tweet_id = ?", tweet.ID).Scan(&tweet_id)

	if err != nil {
		fmt.Println("twitter.go err: " + err.Error())
	}

	return (len(tweet_id) > 0 && err != nil)
}

func writeTweetToDb(db *sql.DB, tweet twitter.Tweet) bool {
	_, err := db.Exec("INSERT INTO `twitter` (text, tweet_id) VALUES (?, ?)", tweet.Text, tweet.ID)

	return (err == nil)
}
