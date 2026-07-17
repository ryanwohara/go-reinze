package runescape

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"unicode/utf8"
)

func Test_WriteNewsToDbWithEmptyNews(t *testing.T) {
	// A failed fetch makes getNews return an empty slice; writing it must
	// be a no-op, not a panic.
	writeNewsToDb(nil, []string{})
}

func Test_GetNewsUnreachableHost(t *testing.T) {
	news := getNews("http://127.0.0.1:1/", "h4 a")

	if len(news) != 0 {
		t.Errorf("Expecting empty slice, received %v", news)
	}
}

func Test_GetNewsNoArticle(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "<html><body><p>maintenance page, no article element</p></body></html>")
	}))
	defer srv.Close()

	news := getNews(srv.URL, "h4 a")

	if len(news) != 0 {
		t.Errorf("Expecting empty slice, received %v", news)
	}
}

func Test_ClampNewsTruncatesToColumnLimits(t *testing.T) {
	news := clampNews([]string{strings.Repeat("a", 300), strings.Repeat("b", 200), "hash1", "oldschool"})

	if utf8.RuneCountInString(news[0]) != 255 {
		t.Errorf("Expecting title truncated to 255 characters, received %d", utf8.RuneCountInString(news[0]))
	}

	if utf8.RuneCountInString(news[1]) != 125 {
		t.Errorf("Expecting url truncated to 125 characters, received %d", utf8.RuneCountInString(news[1]))
	}

	if news[2] != "hash1" || news[3] != "oldschool" {
		t.Errorf("Expecting hash and version untouched, received %v", news[2:])
	}
}

func Test_ClampNewsKeepsShortFields(t *testing.T) {
	news := clampNews([]string{"TWIR - 100", "https://secure.runescape.com/m=news/twir", "hash1", "runescape3"})

	if news[0] != "TWIR - 100" || news[1] != "https://secure.runescape.com/m=news/twir" {
		t.Errorf("Expecting short fields unchanged, received %v", news)
	}
}

func Test_GetNewsNon200(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer srv.Close()

	news := getNews(srv.URL, "h4 a")

	if len(news) != 0 {
		t.Errorf("Expecting empty slice, received %v", news)
	}
}
