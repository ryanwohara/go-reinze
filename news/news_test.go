package news

import (
	"strings"
	"testing"
	"unicode/utf8"

	rss "github.com/mmcdole/gofeed"
)

func Test_NewsFromItemKeepsShortFields(t *testing.T) {
	item := &rss.Item{Title: "Hello World", Link: "https://example.com/story?utm_source=rss"}

	news := newsFromItem(item)

	if news.Title != "Hello World" {
		t.Errorf("Expecting %q, received %q", "Hello World", news.Title)
	}

	if news.Url != "https://example.com/story" {
		t.Errorf("Expecting %q, received %q", "https://example.com/story", news.Url)
	}

	if len(news.Hash) != 5 {
		t.Errorf("Expecting a 5 character hash, received %q", news.Hash)
	}
}

func Test_NewsFromItemFallsBackToGuid(t *testing.T) {
	item := &rss.Item{Title: "Hello World", GUID: "https://example.com/guid?x=1"}

	news := newsFromItem(item)

	if news.Url != "https://example.com/guid" {
		t.Errorf("Expecting %q, received %q", "https://example.com/guid", news.Url)
	}
}

func Test_NewsFromItemTruncatesLongTitle(t *testing.T) {
	item := &rss.Item{Title: strings.Repeat("a", 300), Link: "https://example.com/story"}

	news := newsFromItem(item)

	if utf8.RuneCountInString(news.Title) != 255 {
		t.Errorf("Expecting title truncated to 255 characters, received %d", utf8.RuneCountInString(news.Title))
	}
}

func Test_NewsFromItemTruncatesTitleByRunesNotBytes(t *testing.T) {
	item := &rss.Item{Title: strings.Repeat("é", 300), Link: "https://example.com/story"}

	news := newsFromItem(item)

	if utf8.RuneCountInString(news.Title) != 255 {
		t.Errorf("Expecting title truncated to 255 characters, received %d", utf8.RuneCountInString(news.Title))
	}

	if !utf8.ValidString(news.Title) {
		t.Errorf("Expecting valid UTF-8 after truncation, received %q", news.Title)
	}
}

func Test_NewsFromItemTruncatesLongUrl(t *testing.T) {
	item := &rss.Item{Title: "Hello World", Link: "https://example.com/" + strings.Repeat("a", 300)}

	news := newsFromItem(item)

	if utf8.RuneCountInString(news.Url) != 255 {
		t.Errorf("Expecting url truncated to 255 characters, received %d", utf8.RuneCountInString(news.Url))
	}
}
